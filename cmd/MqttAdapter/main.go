package main

import (
	"encoding/json"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hb9fxq/flexlib-go/obj"
	"net"
	"time"
)

type AppContext struct {
	radioAddr                  string
	myPort                     string
	daxIqChan                  string
	sampleRate                 string
	RadioReponseStreamSequence int
	forwardConnection          net.Conn
	mqttClient                 mqtt.Client
	mqttBroker                 string
	mqttClientId               string
	mqttTopic                  string
}

const NDEF_STRING string = "NDEF"

func publishRaw(appContext *AppContext, message string) {
	rtoken := appContext.mqttClient.Publish(appContext.mqttTopic+"/raw", 0, false, message)
	rtoken.Wait()
}

func main() {

	appContext := new(AppContext)

	flag.StringVar(&appContext.radioAddr, "RADIO", "", "IP ADDRESS OF THE RADIO e.g 192.168.41.8")
	flag.StringVar(&appContext.mqttBroker, "MQTTBROKER", NDEF_STRING, "MQTT Broker conn str.")
	flag.StringVar(&appContext.mqttTopic, "MQTTTOPIC", NDEF_STRING, "MQTT Broker conn str.")
	flag.StringVar(&appContext.mqttClientId, "MQTTCLIENTID", NDEF_STRING, "MQTT Broker conn str.")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(appContext.mqttBroker).SetClientID(appContext.mqttClientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetCleanSession(false)

	appContext.mqttClient = mqtt.NewClient(opts)
	if token := appContext.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	radioContext := new(obj.RadioContext)
	radioContext.RadioAddr = appContext.radioAddr
	radioContext.MyUdpEndpointPort = appContext.myPort
	radioContext.ChannelRadioResponse = make(chan string)
	radioContext.Debug = false
	radioContext.ManualSubscribe = true

	go obj.InitRadioContext(radioContext)

	time.Sleep(2 * time.Second)

	go func(ctx *obj.RadioContext) {
		for {
			response := <-ctx.ChannelRadioResponse
			fmt.Println("F:" + response)

			go publishRaw(appContext, response)
		}
	}(radioContext)

	obj.SendRadioCommand(radioContext, "sub pan all")
	obj.SendRadioCommand(radioContext, "sub slice all")
	obj.SendRadioCommand(radioContext, "sub tx all")

	for {
		time.Sleep(1 * time.Second)

		jsonSlices := make(map[string]interface{})
		radioContext.Slices.Range(func(k interface{}, value interface{}) bool {
			jsonSlices[k.(string)] = value
			return true
		})
		j, err := json.Marshal(&jsonSlices)

		if err != nil {
			fmt.Println(err)
			continue
		}

		stoken := appContext.mqttClient.Publish(appContext.mqttTopic+"/slices", 0, false, j)
		stoken.Wait()

		jsonPanadapters := make(map[string]interface{})
		radioContext.Panadapters.Range(func(k interface{}, value interface{}) bool {
			jsonPanadapters[k.(string)] = value
			return true
		})
		j, err = json.Marshal(&jsonPanadapters)

		if err != nil {
			fmt.Println(err)
			continue
		}

		ptoken := appContext.mqttClient.Publish(appContext.mqttTopic+"/panadapters", 0, false, j)
		ptoken.Wait()

	}

	forever := make(chan bool)
	forever <- true
}
