package pwm_fan

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	viam_utils "go.viam.com/utils"

	"github.com/viam-soleng/viam-fan-controller/utils"
)

var Model = resource.NewModel("viam-soleng", "fan", "pwm")
var PrettyName = "Raspberry Pi Clock Sensor"
var Description = "Simple PWM fan controller for Viam"
var Version = utils.Version

type Config struct {
	resource.Named
	mu               sync.RWMutex
	logger           logging.Logger
	cancelCtx        context.Context
	cancelFunc       func()
	monitor          func()
	done             chan bool
	wg               sync.WaitGroup
	FanPin           board.GPIOPin
	Board            *board.Board
	TemperatureTable map[float64]float64
	Temps            []float64
	Sensor           sensor.Sensor
	SensorValueField string
	SensorValueRegex *regexp.Regexp
}

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *CloudConfig]{Constructor: NewSensor})
}

func NewSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	logger.Infof("Starting %s %s", PrettyName, Version)
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	b := Config{
		Named:      conf.ResourceName().AsNamed(),
		logger:     logger,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		mu:         sync.RWMutex{},
		done:       make(chan bool),
	}

	if err := b.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return &b, nil
}

func (c *Config) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger.Debugf("Reconfiguring %s", PrettyName)

	// In case the module has changed name
	c.Named = conf.ResourceName().AsNamed()

	newConf, err := resource.NativeConfig[*CloudConfig](conf)
	if err != nil {
		return err
	}

	untypedBoard, err := deps.Lookup(resource.NewName(board.API, newConf.BoardName))
	if err != nil {
		c.logger.Errorf("Error looking up board: %s", err)
		return err
	}

	board := untypedBoard.(board.Board)
	fanPin, err := board.GPIOPinByName(newConf.FanPin)
	if err != nil {
		c.logger.Errorf("Error looking up fan pin: %s", err)
		return err
	}

	untypedSensor, err := deps.Lookup(resource.NewName(sensor.API, newConf.SensorName))
	if err != nil {
		c.logger.Errorf("Error looking up sensor: %s", err)
		return err
	}
	sensor := untypedSensor.(sensor.Sensor)

	c.Named = conf.ResourceName().AsNamed()
	c.Board = &board
	c.FanPin = fanPin
	c.Sensor = sensor
	c.SensorValueField = newConf.SensorValueField
	c.SensorValueRegex = regexp.MustCompile(newConf.SensorValueRegex)

	tempTable := make(map[float64]float64)
	temps := make([]float64, 0, len(newConf.TemperatureTable))
	for ts, speed := range newConf.TemperatureTable {
		temp, err := strconv.ParseFloat(ts, 64)
		if err != nil {
			c.logger.Errorf("Error parsing temperature: %s", err)
			return err
		}
		if speed > 1 {
			speed = speed / float64(100)
		}
		tempTable[temp] = speed
		temps = append(temps, temp)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(temps)))

	c.Temps = temps
	c.TemperatureTable = tempTable
	if fanPin.SetPWMFreq(ctx, 1000, nil) != nil {
		c.logger.Errorf("Error setting PWM frequency: %s", err)
		return err
	}

	if c.monitor == nil {
		c.monitor = func() {
			ctx := context.Background()
			c.wg.Add(1)
			defer c.wg.Done()
			for {
				select {
				case <-c.done:
					return
				default:
					readings, err := c.Sensor.Readings(ctx, nil)
					if err != nil {
						c.logger.Errorf("Error getting readings from sensor: %s", err)
						break
					}

					currentTemp, err := utils.ParseCurrentTemperatureFromReadings(ctx, readings, c.SensorValueField, c.SensorValueRegex, c.logger)
					if err != nil {
						c.logger.Errorf("Error parsing current temperature: %s", err)
						break
					}

					desiredSpeed, err := getDesiredSpeed(currentTemp, c.Temps, c.TemperatureTable)
					if err != nil {
						c.logger.Errorf("Error getting desired speed: %s", err)
						break
					}

					c.logger.Debugf("Current temperature: %f, desired speed: %f", currentTemp, desiredSpeed)
					err = c.FanPin.SetPWM(ctx, desiredSpeed, nil)
					if err != nil {
						c.logger.Errorf("Error setting fan speed: %s", err)
					}
				}

				select {
				case <-time.After(100 * time.Millisecond):
					continue
				case <-c.done:
					return
				}
			}
		}

		viam_utils.PanicCapturingGo(c.monitor)
	}

	return nil
}

func (c *Config) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	readings, err := c.Sensor.Readings(ctx, nil)
	if err != nil {
		c.logger.Errorf("Error getting readings from sensor: %s", err)
		return nil, err
	}

	currentTemp, err := utils.ParseCurrentTemperatureFromReadings(ctx, readings, c.SensorValueField, c.SensorValueRegex, c.logger)
	if err != nil {
		c.logger.Errorf("Error parsing current temperature: %s", err)
		return nil, err
	}

	fan_speed, err := c.FanPin.PWM(ctx, nil)
	if err != nil {
		c.logger.Errorf("Error getting fan speed: %s", err)
		return nil, err
	}

	return map[string]interface{}{
		"temperature":   currentTemp,
		"fan_speed_pct": fan_speed * 100,
	}, nil
}

func (c *Config) Close(ctx context.Context) error {
	c.logger.Infof("Shutting down %s", PrettyName)
	c.done <- true
	c.logger.Infof("Notifying monitor to shut down")
	c.wg.Wait()
	c.logger.Info("Monitor shut down")
	return nil
}

func (c *Config) Ready(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return false, nil
}

func getDesiredSpeed(currentTemp float64, temps []float64, tempTable map[float64]float64) (float64, error) {
	for _, targetTemp := range temps {
		if currentTemp >= targetTemp {
			return tempTable[targetTemp], nil
		}
	}

	return 0, errors.New("temperature not found in table")
}
