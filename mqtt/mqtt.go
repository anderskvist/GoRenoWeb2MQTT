package mqtt

import (
	log "github.com/anderskvist/DVIEnergiSmartControl/log"
	renoweb "github.com/anderskvist/GoRenoWeb2MQTT/renoweb"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"

	"fmt"
	"net/url"
	"time"
)

var pubConnection mqtt.Client

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)
	return opts
}

// SendToMQTT will send RenoWeb data to MQTT
func SendToMQTT(cfg *ini.File, pickupPlans renoweb.PickupPlan) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if pubConnection == nil {
		pubConnection = connect("pub", uri)
	}
	for i, pickupPlan := range pickupPlans.List {
		pubConnection.Publish(fmt.Sprintf("renoweb/pickup/%d/ordningnavn", i), 0, false, fmt.Sprintf("%s", pickupPlan.OrdningNavn))
		pubConnection.Publish(fmt.Sprintf("renoweb/pickup/%d/toemningsdato", i), 0, false, fmt.Sprintf("%s", pickupPlan.ToemningsDato))

		t, _ := pickupPlan.ParseToemningsDato()
		pubConnection.Publish(fmt.Sprintf("renoweb/pickup/%d/time", i), 0, false, fmt.Sprintf("%s", t))
		pubConnection.Publish(fmt.Sprintf("renoweb/pickup/%d/hours", i), 0, false, fmt.Sprintf("%.0f", time.Until(t).Hours()))
		pubConnection.Publish(fmt.Sprintf("renoweb/pickup/%d/days", i), 0, false, fmt.Sprintf("%.0f", time.Until(t).Hours()/24))
	}
	pubConnection.Publish("renoweb/pickup/lastupdate", 0, false, time.Now().Format("2006-01-02 15:04:05"))
}
