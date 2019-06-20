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

func pub(t string, s string) {
	if token := pubConnection.Publish(t, 0, false, s); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// SendToMQTT will send RenoWeb data to MQTT
func SendToMQTT(cfg *ini.File, pickupPlans renoweb.PickupPlan) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if pubConnection == nil {
		pubConnection = connect("RenoWeb", uri)
		log.Debug("Connecting to MQTT")
	}
	for i, pickupPlan := range pickupPlans.List {
		pub(fmt.Sprintf("renoweb/pickup/%d/ordningnavn", i), fmt.Sprintf("%s", pickupPlan.OrdningNavn))
		pub(fmt.Sprintf("renoweb/pickup/%d/toemningsdato", i), fmt.Sprintf("%s", pickupPlan.ToemningsDato))

		t, _ := pickupPlan.ParseToemningsDato()
		pub(fmt.Sprintf("renoweb/pickup/%d/time", i), fmt.Sprintf("%s", t))
		pub(fmt.Sprintf("renoweb/pickup/%d/hours", i), fmt.Sprintf("%.0f", time.Until(t).Hours()))
		pub(fmt.Sprintf("renoweb/pickup/%d/days", i), fmt.Sprintf("%.0f", time.Until(t).Hours()/24))
	}
	pub("renoweb/pickup/lastupdate", time.Now().Format("2006-01-02 15:04:05"))
}
