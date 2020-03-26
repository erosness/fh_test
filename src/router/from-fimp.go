package router

import (
	"fmt"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"github.com/thingsplex/thingsplex_service_template/model"
	"strings"
)

type FromFimpRouter struct {
	inboundMsgCh fimpgo.MessageCh
	mqt          *fimpgo.MqttTransport
	instanceId   string
	appLifecycle *model.Lifecycle
	configs      *model.Configs
}

func NewFromFimpRouter(mqt *fimpgo.MqttTransport,appLifecycle *model.Lifecycle,configs *model.Configs) *FromFimpRouter {
	fc := FromFimpRouter{inboundMsgCh: make(fimpgo.MessageCh,5),mqt:mqt,appLifecycle:appLifecycle,configs:configs}
	fc.mqt.RegisterChannel("ch1",fc.inboundMsgCh)
	return &fc
}

func (fc *FromFimpRouter) Start() {

	// TODO: Choose either adapter or app topic

	// ------ Adapter topics ---------------------------------------------
	fc.mqt.Subscribe(fmt.Sprintf("pt:j1/+/rt:dev/rn:%s/ad:1/#",model.ServiceName))
	fc.mqt.Subscribe(fmt.Sprintf("pt:j1/+/rt:ad/rn:%s/ad:1",model.ServiceName))

	// ------ Application topic -------------------------------------------
	//fc.mqt.Subscribe(fmt.Sprintf("pt:j1/+/rt:app/rn:%s/ad:1",model.ServiceName))

	go func(msgChan fimpgo.MessageCh) {
		for  {
			select {
			case newMsg :=<- msgChan:
				fc.routeFimpMessage(newMsg)
			}
		}
	}(fc.inboundMsgCh)
}

func (fc *FromFimpRouter) routeFimpMessage(newMsg *fimpgo.Message) {
	log.Debug("New fimp msg")
	addr := strings.Replace(newMsg.Addr.ServiceAddress,"_0","",1)
	switch newMsg.Payload.Service {
	case "out_lvl_switch" :
		addr = strings.Replace(addr,"l","",1)
		switch newMsg.Payload.Type {
		case "cmd.binary.set":
			// TODO: This is example . Add your logic here or remove
		case "cmd.lvl.set":
			// TODO: This is an example . Add your logic here or remove
		}
	case "out_bin_switch":
		log.Debug("Sending switch")
		// TODO: This is an example . Add your logic here or remove
	case model.ServiceName:
		adr := &fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeAdapter, ResourceName: model.ServiceName, ResourceAddress:"1"}
		switch newMsg.Payload.Type {
		case "cmd.auth.login":
			authReq := model.Login{}
			err := newMsg.Payload.GetObjectValue(&authReq)
			if err != nil {
				log.Error("Incorrect login message ")
				return
			}
			status := model.AuthStatus{
				Status:    "AUTHENTICATED",
				ErrorText: "",
				ErrorCode: "",
			}
			if authReq.Username != "" && authReq.Password != ""{
				// TODO: This is an example . Add your logic here or remove
			}else {
				status.Status = "ERROR"
				status.ErrorText = "Empty username or password"
			}

			msg := fimpgo.NewMessage("evt.auth.status_report",model.ServiceName,fimpgo.VTypeObject,status,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				// if response topic is not set , sending back to default application event topic
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.app.get_manifest":
			mode,err := newMsg.Payload.GetStringValue()
			if err != nil {
				log.Error("Incorrect request format ")
				return
			}
			manifest := model.NewManifest()
			manifest.LoadFromFile("../testdata/app-manifest.json")
			if mode == "manifest_and_states" {
				manifest.AppState = *fc.appLifecycle.GetAllStates()
			}
			msg := fimpgo.NewMessage("evt.app.manifest_report",model.ServiceName,fimpgo.VTypeObject,manifest,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				// if response topic is not set , sending back to default application event topic
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.app.get_state":
			msg := fimpgo.NewMessage("evt.app.manifest_report",model.ServiceName,fimpgo.VTypeObject,fc.appLifecycle.GetAllStates(),nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				// if response topic is not set , sending back to default application event topic
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.config.get_extended_report":

			msg := fimpgo.NewMessage("evt.config.extended_report",model.ServiceName,fimpgo.VTypeObject,fc.configs,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.config.extended_set":
			err :=newMsg.Payload.GetObjectValue(fc.configs)
			if err != nil {
				// TODO: This is an example . Add your logic here or remove
				log.Error("Can't parse configuration object")
				return
			}
			fc.configs.SaveToFile()
			log.Debugf("App reconfigured . New parameters : %v",fc.configs)
			// TODO: This is an example . Add your logic here or remove
			configReport := model.ConfigReport{
				OpStatus: "OK",
				AppState:  *fc.appLifecycle.GetAllStates(),
			}
			msg := fimpgo.NewMessage("evt.app.config_report",model.ServiceName,fimpgo.VTypeObject,configReport,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.log.set_level":
			// Configure log level
			level , err :=newMsg.Payload.GetStringValue()
			if err != nil {
				return
			}
			logLevel, err := log.ParseLevel(level)
			if err == nil {
				log.SetLevel(logLevel)
				fc.configs.LogLevel = level
				fc.configs.SaveToFile()
			}
			log.Info("Log level updated to = ",logLevel)

		case "cmd.system.connect":
			// This is optional operation.
			var errStr string
			status := "ok"
			fc.appLifecycle.PublishEvent(model.EventConfigured,"from-fimp-router",nil)
			val := map[string]string{"status":status,"error":errStr}
			msg := fimpgo.NewStrMapMessage("evt.system.connect_report",model.ServiceName,val,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.network.get_all_nodes":
			// TODO: This is an example . Add your logic here or remove
		case "cmd.thing.get_inclusion_report":
			//nodeId , _ := newMsg.Payload.GetStringValue()
			// TODO: This is an example . Add your logic here or remove
		case "cmd.thing.inclusion":
			//flag , _ := newMsg.Payload.GetBoolValue()
			// TODO: This is an example . Add your logic here or remove
		case "cmd.thing.delete":
			// remove device from network
			val,err := newMsg.Payload.GetStrMapValue()
			if err != nil {
				log.Error("Wrong msg format")
				return
			}
			deviceId , ok := val["address"]
			if ok {
				// TODO: This is an example . Add your logic here or remove
				log.Info(deviceId)
			}else {
				log.Error("Incorrect address")

			}
		}

	}

}


