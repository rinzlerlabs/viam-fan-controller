package main

import (
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	"github.com/rinzlerlabs/viam-fan-controller/on_off_fan"
	"github.com/rinzlerlabs/viam-fan-controller/pwm_fan"

	raspiutils "github.com/rinzlerlabs/viam-fan-controller/utils"
	moduleutils "github.com/thegreatco/viamutils/module"
)

func main() {
	logger := module.NewLoggerFromArgs(raspiutils.LoggerName)
	logger.Infof("Starting RinzlerLabs Fan Controller Module %v", raspiutils.Version)
	moduleutils.AddModularResource(on_off_fan.API, on_off_fan.Model)
	moduleutils.AddModularResource(pwm_fan.API, pwm_fan.Model)
	utils.ContextualMain(moduleutils.RunModule, logger)
}
