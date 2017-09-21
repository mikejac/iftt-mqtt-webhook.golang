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
	 "github.com/twinj/uuid"
 )
 
type MsgbusStatus int

//
//
type MqttOptions struct {
	Server 				string					// name or ip of the MQTT server
	Port 				int						// portnumber of the MQTT server (usually 1883)
	ClientId			string					// MQTT client id
	Keepalive 			int						// MQTT keep-alive interval in seconds

	Domain 		    	string					// Very first part of all MQTT topics
	Nodename 			string					// Our nodename
	StatusInterval		int

	StateChangeCallback	StateChangeCallback
	NodeChangeCallback	NodeChangeCallback
}

//
//
func NewMqttOptions() *MqttOptions {
	o := &MqttOptions{
		Server:			"localhost",
		Port:			1883,
		ClientId:		"",
		Keepalive:		60,
		Domain:			"domain",
		Nodename:		uuid.NewV4().String(),
		StatusInterval:	60,
	}

	return o
}
 
//
func (o *MqttOptions) SetServer(server string) (*MqttOptions) {
	o.Server = server
	return o
}
 
//
func (o *MqttOptions) SetPort(port int) (*MqttOptions) {
	o.Port = port
	return o
}
 
//
func (o *MqttOptions) SetClientId(clientId string) (*MqttOptions) {
	o.ClientId = clientId
	return o
}
 
//
func (o *MqttOptions) SetKeepalive(keepalive int) (*MqttOptions) {
	o.Keepalive = keepalive
	return o
}
 
//
func (o *MqttOptions) SetDomain(domain string) (*MqttOptions) {
	o.Domain = domain
	return o
}

//
func (o *MqttOptions) SetStatusInterval(interval int) (*MqttOptions) {
	o.StatusInterval = interval
	return o
}

//
func (o *MqttOptions) SetNodename(nodename string) (*MqttOptions) {
	o.Nodename = nodename
	return o
}
 
type StateChangeCallback func(connected bool)

//
func (o *MqttOptions) SetStateChangeCallback(fn StateChangeCallback) (*MqttOptions) {
	o.StateChangeCallback = fn
	return o
}
 
type NodeChangeCallback func(nodename string, status MsgbusStatus, uptime int64)
 
//
func (o *MqttOptions) SetNodeChangeCallback(fn NodeChangeCallback) (*MqttOptions) {
	o.NodeChangeCallback = fn
	return o
}
