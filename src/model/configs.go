package model

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/thingsplex/thingsplex_service_template/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const ServiceName  = "thingsplex_service_template"

type Configs struct {
	path                  string
	InstanceAddress       string `json:"instance_address"`
	MqttServerURI         string `json:"mqtt_server_uri"`
	MqttUsername          string `json:"mqtt_server_username"`
	MqttPassword          string `json:"mqtt_server_password"`
	MqttClientIdPrefix    string `json:"mqtt_client_id_prefix"`
	LogFile               string `json:"log_file"`
	LogLevel              string `json:"log_level"`
	LogFormat             string `json:"log_format"`
	WorkDir               string `json:"work_dir"`
	ConfiguredAt          string `json:"configured_at"`
	ConfiguredBy          string `json:"configured_by"`
	Param1                bool   `json:"param_1"`
	Param2                string `json:"param_2"`
}

func NewConfigs(path string) *Configs {
	return &Configs{path:path}
}

func (cf * Configs) LoadFromFile() error {
	configFileBody, err := ioutil.ReadFile(cf.path)
	if err != nil {
		cf.InitDefault()
		return cf.SaveToFile()
	}
	err = json.Unmarshal(configFileBody, cf)
	if err != nil {
		return err
	}
	return nil
}

func (cf *Configs) SaveToFile() error {
	cf.ConfiguredBy = "auto"
	cf.ConfiguredAt = time.Now().Format(time.RFC3339)
	bpayload, err := json.Marshal(cf)
	err = ioutil.WriteFile(cf.path, bpayload, 0664)
	if err != nil {
		return err
	}
	return err
}

func (cf *Configs) GetDataDir()string {
	return filepath.Join(cf.WorkDir,"data")
}

func (cf *Configs) GetDefaultDir()string {
	return filepath.Join(cf.WorkDir,"defaults")
}

func (cf * Configs) LoadDefaults()error {
	configFile := filepath.Join(cf.WorkDir,"data","config.json")
	os.Remove(configFile)
	log.Info("Config file doesn't exist.Loading default config")
	defaultConfigFile := filepath.Join(cf.WorkDir,"defaults","config.json")
	return utils.CopyFile(defaultConfigFile,configFile)
}

func (cf *Configs) InitDefault() {
	cf.InstanceAddress = "1"
	cf.MqttServerURI = "tcp://localhost:1883"
	cf.MqttClientIdPrefix = "thingsplex_service_template"
	cf.LogFile = "/var/log/thingsplex/thingsplex_service_template/thingsplex_service_template.log"
	cf.WorkDir = "/opt/thingsplex/thingsplex_service_template"
	cf.LogLevel = "debug"
	cf.LogFormat = "text"
	cf.Param1 = true
	cf.Param2 = "test"
}

func (cf *Configs) IsConfigured()bool {
	// TODO : Add logic here
	return true
}

type ConfigReport struct {
	OpStatus string `json:"op_status"`
	AppState AppStates `json:"app_state"`
}