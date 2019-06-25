package main

import (
	"os"
	"time"

	"github.com/anderskvist/GoHelpers/version"

	"github.com/anderskvist/GoRenoWeb2MQTT/mqtt"
	"github.com/anderskvist/GoRenoWeb2MQTT/renoweb"

	"github.com/anderskvist/DVIEnergiSmartControl/log"
	ini "gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load(os.Args[1])

	if err != nil {
		log.Criticalf("Fail to read file: %v", err)
		os.Exit(1)
	}

	log.Infof("GoRenoWeb version: %s.\n", version.Version)

	poll := cfg.Section("main").Key("poll").MustInt(60)
	log.Infof("Polltime is %d seconds.\n", poll)

	reno, err := renoweb.NewClient(
		renoweb.SetDebugLogger(log.Debugf),
		renoweb.SetHostname(cfg.Section("renoweb").Key("hostname").String()),
	)
	if err != nil {
		log.Criticalf("Unable to instantiate renoweb client: %s", err.Error())
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(poll) * time.Second)
	for ; true; <-ticker.C {
		log.Notice("Tick")
		log.Info("Getting data from RenoWeb")
		addressID, _ := reno.AddressID(cfg.Section("renoweb").Key("address").String())
		pickupPlans, _ := reno.PickupPlan(addressID)
		log.Info("Done getting data from RenoWeb")

		log.Info("Sending data to MQTT")
		mqtt.SendToMQTT(cfg, *pickupPlans)
		log.Info("Done sending to MQTT")
	}
}
