package main

import (
	"github.com/Shopify/goose/logger"
	"github.com/Shopify/goose/safely"

	"github.com/CovidShield/server/pkg/app"
	"github.com/CovidShield/server/pkg/telemetry"
)

var log = logger.New("main")

func main() {
	defer safely.Recover() // panics -> bugsnag

	log(nil, nil).Info("starting")

	cleanupTracer := telemetry.InitTracer()
	defer cleanupTracer()

	cleanupMeter := telemetry.InitMeter()
	defer cleanupMeter()

	mainApp := app.NewBuilder().WithSubmission().Build()

	err := mainApp.RunAndWait()
	defer log(nil, err).Info("final message before shutdown")
}
