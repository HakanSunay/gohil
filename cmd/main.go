package main

import (
	"context"

	"github.com/HakanSunay/gohil/logger"
)

var ctx = context.Background()

func init() {
	logger.InitGlobalLogging()
	ctx = logger.PutLoggerInContext(ctx)
}

func main() {
	log := logger.GetFromContext(ctx)
	log.Infof("Starting gohil...")
	log.Infof("Terminating gohil...")
	return
}
