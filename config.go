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
	"github.com/go-ini/ini"
	"github.com/mikejac/log.golang"
)

//
//
func NewConfig() (config *DispatcherConfiguration) {
    config					= &DispatcherConfiguration{}
    config.MqttOptions	    = NewMqttOptions()
    config.httpIp 			= ""
    config.httpPort 		= "8080"
    config.useTLS           = false
    
    return config
}

//
//
func (config *DispatcherConfiguration) ReadConfig(configfile string) (err error) {
	cfg, err := ini.Load(configfile)
	if err != nil {
    	log.Infof("error: %s\n", err.Error())
		return err
	}
    
	/******************************************************************************************************************
	 * MQTT settings
	 *
	 */    
    if cfg.Section("mqtt").HasKey("clientid") {
        config.MqttOptions.SetClientId(cfg.Section("mqtt").Key("clientid").String())
    }

    if cfg.Section("mqtt").HasKey("server") {
        config.MqttOptions.SetServer(cfg.Section("mqtt").Key("server").String())
    }

    if cfg.Section("mqtt").HasKey("port") {
        port, _ := cfg.Section("mqtt").Key("port").Int()
        config.MqttOptions.SetPort(port)
    }

    if cfg.Section("mqtt").HasKey("keepalive") {
        keepalive, _ := cfg.Section("mqtt").Key("keepalive").Int()
        config.MqttOptions.SetKeepalive(keepalive)
    }

	/******************************************************************************************************************
	 * MsgBus settings
	 *
	 */    
     if cfg.Section("msgbus").HasKey("nodename") {
        config.MqttOptions.SetNodename(cfg.Section("msgbus").Key("nodename").String())
    }

    if cfg.Section("msgbus").HasKey("domain") {
        config.MqttOptions.SetDomain(cfg.Section("msgbus").Key("domain").String())
    }

    if cfg.Section("msgbus").HasKey("status_interval") {
        interval, _ := cfg.Section("msgbus").Key("status_interval").Int()
        config.MqttOptions.SetStatusInterval(interval)
    }

	/******************************************************************************************************************
	 * HTTP settings
	 *
	 */
    if cfg.Section("http").HasKey("addr") {
        config.httpIp = cfg.Section("http").Key("addr").String()
    }

    if cfg.Section("http").HasKey("port") {
        config.httpPort = cfg.Section("http").Key("port").String()
    }

    if cfg.Section("http").HasKey("use_tls") {
        useTLS, _ := cfg.Section("http").Key("use_tls").Bool()
        config.useTLS = useTLS
    }

	/******************************************************************************************************************
	 * TLS settings
	 *
     */
    if config.useTLS {
        if cfg.Section("tls").HasKey("cert") {
            config.certFile = cfg.Section("tls").Key("cert").String()
        }

        if cfg.Section("tls").HasKey("key") {
            config.keyFile = cfg.Section("tls").Key("key").String()
        }
    }

	/******************************************************************************************************************
	 * API Keys
	 *
     */
    names := cfg.Section("apikeys").KeyStrings()
    
    for _, n := range names {
        apikey := cfg.Section("apikeys").Key(n).String()
        
        config.apikeys = append(config.apikeys, apikey)
    }
 
    return nil
}
