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
	"github.com/mikejac/log.golang"
)

type DispatcherConfiguration struct {
	MqttOptions			*MqttOptions
	
	httpIp				string
	httpPort			string
	
	useTLS				bool
	certFile			string
	keyFile				string

	apikeys				[]string
}

type Dispatcher struct {
    config 				*DispatcherConfiguration
    mqtt				*Mqtt
	
    exit 				chan bool
    
	chanMqttStateChange chan bool
	chanMqttNodeChange 	chan bool

	httpLocation		chan Location
}

//
//
func NewDispatcher(config *DispatcherConfiguration, exit chan bool) (dispatcher *Dispatcher) {
	log.Debugf("NewDispatcher(): begin")

    dispatcher = &Dispatcher{config: config, exit: exit}
	
	dispatcher.httpLocation  = make(chan Location)

	// set callbacks
	dispatcher.config.MqttOptions.SetStateChangeCallback(dispatcher.stateChangeCallback)
	dispatcher.config.MqttOptions.SetNodeChangeCallback(dispatcher.nodeChangeCallback)

	log.Debugf("NewDispatcher(): end")
	
	return dispatcher
}
//
//
func (dispatcher *Dispatcher) Run() (err error) {
	log.Debugf("DispatcherData::Run(): begin")

    dispatcher.mqtt, err = NewConnector(dispatcher.config.MqttOptions)
	if err != nil {
		log.Infof("Dispatcher::Run(): NewConnector() error ", err.Error())
		return err
	}

	dispatcher.mqtt.Connect()
	
	httpServer := NewHttpServer(dispatcher.config, dispatcher)
	if httpServer == nil {
		log.Infof("Dispatcher::Run(): NewHttpServer() error ", err.Error())
		return err
	}

	if err := httpServer.Start(); err != nil {
		log.Infof("Dispatcher::Run(): failed to start HTTP server; ", err.Error())
		return err
	}
	
	var shouldRun = true
	
	for shouldRun {
		select {
			/******************************************************************************************************************
			 * incoming http data
			 *
			 */
			case r := <- dispatcher.httpLocation:
				log.Debugf("Dispatcher::Run(): got 'httpLocation'")
				log.Debugf("Dispatcher::Run(): r = %+v", r)

				dispatcher.mqtt.PublishUpdate(r.dataId, r)

			/******************************************************************************************************************
			 * exit
			 *
			 */
			case <- dispatcher.exit:
				log.Debugf("Dispatcher::Run(): exit")
				shouldRun = false
				break
		}
	}	

	log.Debugf("Dispatcher::Run(): end")
    
    return nil
}

/******************************************************************************************************************
* MQTT transitions
*
*/

//
//
func (dispatcher *Dispatcher) stateChangeCallback(connected bool) {
	log.Debugf("Dispatcher::stateChangeCallback()")
}

func (dispatcher *Dispatcher) nodeChangeCallback(nodename string, status MsgbusStatus, uptime int64) {
	log.Debugf("Dispatcher::nodeChangeCallback()")
}
