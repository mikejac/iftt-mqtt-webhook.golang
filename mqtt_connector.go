/*
 * Copyright (c) 2017 Michael Jacobsen (github.com/mikejac)
 *
 * This file is part of esp8266upgrader.golang.
 *
 * iftt-mqtt-webhook.golang is free software: you can redistribute
 * it and/or modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * iftt-mqtt-webhook.golang is distributed in the hope that it will
 * be useful, but WITHOUT ANY WARRANTY; without even the implied warranty
 * of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with esp8266upgrader.golang.  If not,
 * see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"errors"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"github.com/mikejac/log.golang"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//
//
type Mqtt struct {
	// Paho MQTT
	client					MQTT.Client
	options  	 			*MQTT.ClientOptions
	qos 					byte
	
	// MessageBus data
	domain       			string
	nodename         		string

	// other
	statusInterval			int
	startTime    			time.Time

	stateChangeCallback		StateChangeCallback
	nodeChangeCallback		NodeChangeCallback
}

type statusUpdate struct {
	Status	string	`json:"status"`
	Uptime	int64	`json:"uptime"`
}

//
//
func NewConnector(options *MqttOptions) (*Mqtt, error) {
	mqtt := &Mqtt{}

	mqtt.qos					= 1
	mqtt.stateChangeCallback	= options.StateChangeCallback
	mqtt.nodeChangeCallback		= options.NodeChangeCallback
	mqtt.statusInterval			= options.StatusInterval
	mqtt.startTime				= time.Now()
	mqtt.domain					= options.Domain
	mqtt.nodename				= options.Nodename

	var clientId string
	
	if options.ClientId == "" {
		hostname, _ := os.Hostname()
		clientId     = hostname + strconv.Itoa(time.Now().Second())
	} else {
		clientId = options.ClientId
	}
	
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://" + options.Server + ":" + strconv.Itoa(options.Port))
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(time.Duration(options.Keepalive) * time.Second)
	opts.SetDefaultPublishHandler(mqtt.onMessage)
	opts.SetOnConnectHandler(mqtt.onConnect)
	opts.SetConnectionLostHandler(mqtt.onDisconnect)

	mqtt.options = opts
	
	if mqtt.client = MQTT.NewClient(mqtt.options); mqtt.client == nil {
		return nil, errors.New("could not create Paho MQTT client")
	}
	
	return mqtt, nil
}
//
//
func (mqtt *Mqtt) Connect() error {
	if mqtt == nil {
		return errors.New("'mqtt' is nil")
	}

	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	ticker := time.NewTicker(time.Second * time.Duration(mqtt.statusInterval))
    go func() {
        for t := range ticker.C {
			log.Debug("mqtt::Connect(): tick at ", t)

			status := statusUpdate{
				Status:	"online",
				Uptime:	time.Now().Unix() - mqtt.startTime.Unix(),
			}

			mqtt.PublishUpdate("Status", status)
        }
	}()
	
	return nil
}
//
//
func (mqtt *Mqtt) Close() error {
	log.Debugf("mqtt::Close(): begin")
	
	status := statusUpdate{
		Status:	"offline",
		Uptime:	time.Now().Unix() - mqtt.startTime.Unix(),
	}

	mqtt.PublishUpdate("Status", status)
	
	mqtt.client.Disconnect(250)
	
	log.Debugf("mqtt::Close(): end")

	return nil
}
//
//
func (mqtt *Mqtt) PublishUpdate(dataId string, data interface{}) error {
	topic := mqtt.topicUpdate(dataId)

	log.Debugf("mqtt::PublishUpdate(): topic = %s", topic)
	
	b, err := json.Marshal(data)
	if err != nil {
		log.Info("mqtt::PublishUpdate(): marshal error = ", err)
		return err
	}

	log.Debugf("mqtt::PublishUpdate(): b = %s", string(b[:]))
	
	if token := mqtt.client.Publish(topic, mqtt.qos, false, b); token.Wait() && token.Error() != nil {
		log.Debugf("mqtt::PublishUpdate(): err = %s", token.Error().Error())
		return token.Error()
	}

	return nil
}

/******************************************************************************************************************
 * MQTT topics
 *
 */

const (
	msgbusSelf                  string = "msgbus"
	msgbusVersion               string = "v2"

	msgbusDestBroadcast         string = "broadcast"

	msgbusUpdate                string = "Update"
	msgbusWrite                 string = "Write"
	msgbusRead                  string = "Read"
	msgbusRPC                   string = "rpc"
)

func (mqtt *Mqtt) topicUpdate(dataId string) string {
	topic := 	mqtt.domain + "/" + 
				msgbusSelf + "/" +
				msgbusVersion + "/" +
				msgbusDestBroadcast + "/" +
				mqtt.nodename + "/" +
				dataId + "." + msgbusUpdate

	return topic
}

/******************************************************************************************************************
 * MQTT event handlers
 *
 */

//
//
func (mqtt *Mqtt) onConnect(client MQTT.Client) {
	log.Debugf("mqtt::onConnect()")

	if mqtt.stateChangeCallback != nil {
		mqtt.stateChangeCallback(true)
	}
}
//
//
func (mqtt *Mqtt) onDisconnect(client MQTT.Client, err error) {
	log.Debugf("mqtt::onDisonnect()")
	
	if mqtt.stateChangeCallback != nil {
		mqtt.stateChangeCallback(false)
	}
}
func (mqtt *Mqtt) onMessage(client MQTT.Client, msg MQTT.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Info("mqtt::onMessage(): panic recovered; ", r)
		}
	}()

	log.Debugf("mqtt::onMessage()")
}
