package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/satori/go.uuid"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {

	c := newOneClient("192.168.1.240:1883")

	c.Subscribe("equipment/#", 2, func(client mqtt.Client, message mqtt.Message) {
		log.Infof("message:%v", string(message.Payload()))
	})
	//t := time.NewTicker(time.Second * 5)
	//for range t.C {
	//	p := pprof.Lookup("goroutine")
	//
	//	log.Infof("goroutine:%v", p.Count())
	//	p.WriteTo(os.Stdout, 1)
	//}
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
	// connection lost callback
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		log.Infof("error is :%v", e)
		for {
			// client not connectErr mqtt,need connectErr to server
			if !client.IsConnectionOpen() {
				// until connectErr to server
				if connectErr(client) {
					log.Errorf("can not client mqtt server!! pelease check! address = %v", addr)
				} else {

					log.Infof("connectErr to mqtt server")

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
	if connectErr(client) {
		log.Fatalf("can not connectErr mqtt server!! pelease check! address = %s!", addr)
		client.Disconnect(250)
	} else {
		log.Warnf("connectErr to mqtt ok clientid = %s!", clientid)
	}
	return client
}
func connectErr(client mqtt.Client) bool {
	token := client.Connect()
	return token.Wait() && token.Error() != nil
}
