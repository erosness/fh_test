package model

import (
	"encoding/json"
	"github.com/futurehomeno/fimpgo/fimptype"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Manifest struct {
	Configs  []AppConfig   `json:"configs"`
	UIBlocks []AppUBLock   `json:"ui_blocks"`
	Auth     AppAuth       `json:"auth"`
	InitFlow []string      `json:"init_flow"`
	Services []AppServices `json:"services"`
	AppState AppStates     `json:"app_state"`
}

type AppConfig struct {
	ID          string            `json:"id"`
	Label       MultilingualLabel `json:"label"`
	ValT        string            `json:"val_t"`
	UI          AppConfigUI       `json:"ui"`
	Val         Value             `json:"val"`
	IsRequired  bool              `json:"is_required"`
	ConfigPoint string            `json:"config_point"`
}

type MultilingualLabel map[string]string

type AppAuth struct {
	Type         string `json:"type"`
	RedirectURL  string `json:"redirect_url"`
	ClientID     string `json:"client_id"`
	PartnerID    string `json:"partner_id"`
	AuthEndpoint string `json:"auth_endpoint"`
}

type AppServices struct {
	Name       string               `json:"name"`
	Alias      string               `json:"alias"`
	Address    string               `json:"address"`
	Interfaces []fimptype.Interface `json:"interfaces"`
}

type Value struct {
	Default interface{} `json:"default"`
}

type AppConfigUI struct {
	Type   string      `json:"type"`
	Select interface{} `json:"select"`
}

type AppUBLock struct {
	Header  MultilingualLabel `json:"label"`
	Text    MultilingualLabel `json:"text"`
	Configs []string          `json:"configs"`
	Footer  MultilingualLabel `json:"footer"`
}

func NewManifest() *Manifest {
	return &Manifest{}
}

func (m *Manifest) LoadFromFile(filePath string) error {
	log.Debug("<manifest> Loading flow from file : ", filePath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("<manifest> Can't open manifest file.")
		return err
	}
	err = json.Unmarshal(file, m)
	if err != nil {
		log.Error("<FlMan> Can't unmarshal manifest file.")
		return err
	}
	return nil
}

func (m *Manifest)SaveToFile(filePath string) error {
	flowMetaByte,err := json.Marshal(m)
	if err != nil {
		log.Error("<manifest> Can't marshal imported file ")
		return err
	}
	log.Debugf("<manifest> Saving manifest to file %s :", filePath)
	err = ioutil.WriteFile(filePath, flowMetaByte, 0644)
	if err != nil {
		log.Error("<manifest>Can't save flow to file . Error : ", err)
		return err
	}
	return nil
}