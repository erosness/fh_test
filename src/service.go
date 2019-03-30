package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alivinco/thingsplex_service_template/model"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"time"
)

func SetupLog(logfile string, level string, logFormat string) {
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.999"})
	} else {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true, TimestampFormat: "2006-01-02T15:04:05.999"})
	}

	logLevel, err := log.ParseLevel(level)
	if err == nil {
		log.SetLevel(logLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	if logfile != "" {
		l := lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    5, // megabytes
			MaxBackups: 2,
		}
		log.SetOutput(&l)
	}

}

func main() {
	configs := model.Configs{}
	var configFile string
	flag.StringVar(&configFile, "c", "", "Config file")
	flag.Parse()
	if configFile == "" {
		configFile = "./config.json"
	} else {
		fmt.Println("Loading configs from file ", configFile)
	}
	configFileBody, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(configFileBody, &configs)
	if err != nil {
		fmt.Print(err)
		panic("Can't load config file.")
	}

	SetupLog(configs.LogFile, configs.LogLevel, configs.LogFormat)
	log.Info("--------------Starting thingsplex_service_template----------------")

	mqtt := fimpgo.NewMqttTransport(configs.MqttServerURI,configs.MqttClientIdPrefix,configs.MqttUsername,configs.MqttPassword,true,1,1)
	err = mqtt.Start()
	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "thingsplex_service_template", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)
	if err != nil {
		log.Error("Can't connect to broker. Error:",err.Error())
	}else {
		log.Info("Connected")
	}
	mqtt.Stop()
	time.Sleep(10*time.Second)
}
