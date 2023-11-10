package main

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	"github.com/viam-soleng/viam-fan-controller/pwm_fan"
	"github.com/viam-soleng/viam-fan-controller/simple_fan"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("viam-raspi-utils"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	custom_module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = custom_module.AddModelFromRegistry(ctx, sensor.API, simple_fan.Model)
	if err != nil {
		return err
	}

	err = custom_module.AddModelFromRegistry(ctx, sensor.API, pwm_fan.Model)
	if err != nil {
		return err
	}

	err = custom_module.Start(ctx)
	defer custom_module.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
