package main

import (
	"context"
	"os"

	"github.com/HakanSunay/gohil/logger"
	"github.com/HakanSunay/gohil/shell"
)

var ctx = context.Background()

func init() {
	logger.InitGlobalLogging()
	ctx = logger.PutLoggerInContext(ctx)
}

func main() {
	log := logger.GetFromContext(ctx)
	log.Infof("Starting gohil...")
	shell.Start(ctx, os.Stdin, os.Stdout)
	log.Infof("Terminating gohil...")
}
