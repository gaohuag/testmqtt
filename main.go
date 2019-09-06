package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	c := newOneClient("192.168.1.240:1883")
	c.Subscribe("equipment/#", 2, func(client mqtt.Client, message mqtt.Message) {
		log.Infof("message:%v", string(message.Payload()))
	})
	select {}
}

func newOneClient(addr string) mqtt.Client {
	uuid := uuid.NewV4()
	clientid := uuid.String()
	opts := mqtt.NewClientOptions().AddBroker(addr).SetClientID(clientid)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetProtocolVersion(4)
	opts.SetCleanSession(true)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		log.Infof("error is :%v", e)
		for {
			if !client.IsConnectionOpen() {
				if token := client.Connect(); token.Wait() && token.Error() != nil {
					log.Errorf("can not client mqtt server!! pelease check! address = %v", addr)
				} else {

					log.Infof("connect to mqtt server")

					break
				}
			} else {
				log.Infof("mqtt is connected,%v", client.IsConnectionOpen())

				break
			}
			time.Sleep(time.Second * 5)
		}
	})
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("can not connect mqtt server!! pelease check! address = %s!", addr)
		client.Disconnect(250)
	} else {
		log.Warnf("connect to mqtt ok clientid = %s!", clientid)
	}
	return client
}
