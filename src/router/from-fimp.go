package router

import (
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"github.com/thingsplex/thingsplex_service_template/model"
	"strings"
)

const ServiceName  = "thingsplex_service_template"

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
	fc.mqt.Subscribe("pt:j1/+/rt:dev/rn:thingsplex_service_template/ad:1/#")
	fc.mqt.Subscribe("pt:j1/+/rt:ad/rn:cthingsplex_service_template/ad:1")
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
			//val,_ := newMsg.Payload.GetBoolValue()
			// TODO: Add your logic here

			//log.Debug("Status code = ",respH.StatusCode)
		case "cmd.lvl.set":
			//val,_ := newMsg.Payload.GetIntValue()
			// TODO: Add your logic here
		}

	case "out_bin_switch":
		log.Debug("Sending switch")
		//val,_ := newMsg.Payload.GetBoolValue()
		// TODO: Add your logic here
	case ServiceName:
		adr := &fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeAdapter, ResourceName: ServiceName, ResourceAddress:"1"}
		switch newMsg.Payload.Type {
		case "cmd.auth.login":
			reqVal, err := newMsg.Payload.GetStrMapValue()
			status := "ok"
			if err != nil {
				log.Error("Incorrect login message ")
				return
			}
			username,_ := reqVal["username"]
			password,_ := reqVal["password"]
			if username != "" && password != ""{
				// TODO: Add your logic here
			}
			fc.configs.SaveToFile()
			if err != nil {
				status = "error"
			}
			msg := fimpgo.NewStringMessage("evt.system.login_report",ServiceName,status,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}
		case "cmd.system.get_connect_params":
			val := map[string]string{"host":"","username":"","password":""}
			msg := fimpgo.NewStrMapMessage("evt.system.connect_params_report",ServiceName,val,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}
		case "cmd.config.set":
			fallthrough
		case "cmd.system.connect":
			_, err := newMsg.Payload.GetStrMapValue()
			var errStr string
			status := "ok"
			//if err != nil {
			//	log.Error("Incorrect login message ")
			//	errStr = err.Error()
			//}
			//host,_ := reqVal["host"]
			//username,_ := reqVal["username"]
			//password,_ := reqVal["password"]

			//if username != ""{
			fc.appLifecycle.PublishEvent(model.EventConfigured,"from-fimp-router",nil)
			//}
			fc.configs.SaveToFile()
			if err != nil {
				status = "error"
			}
			val := map[string]string{"status":status,"error":errStr}
			msg := fimpgo.NewStrMapMessage("evt.system.connect_report",ServiceName,val,nil,nil,newMsg.Payload)
			if err := fc.mqt.RespondToRequest(newMsg.Payload,msg); err != nil {
				fc.mqt.Publish(adr,msg)
			}

		case "cmd.network.get_all_nodes":
			// TODO: Add your logic here
		case "cmd.thing.get_inclusion_report":
			//nodeId , _ := newMsg.Payload.GetStringValue()
			// TODO: Add your logic here
		case "cmd.thing.inclusion":
			//flag , _ := newMsg.Payload.GetBoolValue()
			// TODO: Add your logic here
		case "cmd.thing.delete":
			// remove device from network
			val,err := newMsg.Payload.GetStrMapValue()
			if err != nil {
				log.Error("Wrong msg format")
				return
			}
			deviceId , ok := val["address"]
			if ok {
				// TODO: Add your logic here
				log.Info(deviceId)
			}else {
				log.Error("Incorrect address")

			}

		}
		//

	}

}


