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
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"	
	"github.com/mikejac/log.golang"
)

type Location struct {
	Who		string	`json:"who"`
	Area	string	`json:"area"`
	Type	string	`json:"type"`

	dataId	string
}

type HttpServerData struct {
    config 		*DispatcherConfiguration
	dispatcher	*Dispatcher

	addr		string
}

type HttpHandler struct {
	
}

//
//
func NewHttpServer(config *DispatcherConfiguration, dispatcher *Dispatcher) (server *HttpServerData) {
	log.Debugf("NewHttpServer(): begin")

	server = &HttpServerData{config: config, dispatcher: dispatcher, addr: config.httpIp + ":" + config.httpPort}

	log.Debugf("NewHttpServer(): addr = '%s'\n", server.addr)
	log.Debugf("NewHttpServer(): end")

	return server
}
//
//
func (server *HttpServerData) Start() (err error) {
	log.Debugf("HttpServerData::Start(): begin")
	
	go func() {
		log.Debugf("HttpServerData::Start(): go func begin")
		
		mux := http.NewServeMux()
		
		mux.Handle("/ifttt/", server)
		
		if server.config.useTLS {
			log.Debugf("HttpServerData::Start(): using TLS")
			if err = http.ListenAndServeTLS(server.addr, server.config.certFile, server.config.keyFile, mux); err != nil {
				log.Info(err.Error())
			}			
		} else {
			log.Debugf("HttpServerData::Start(): using plaintext")
			if err = http.ListenAndServe(server.addr, mux); err != nil {
				log.Info(err.Error())
			}
		}

		log.Debugf("HttpServerData::Start(): go func end")
	}()

	log.Debugf("HttpServerData::Start(): end")

	return err
}
//
//
func (server *HttpServerData) Stop() (err error) {
	return nil
}
//
//
func (server *HttpServerData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rr := recover(); rr != nil {
			log.Info("HttpServerData::ServeHTTP(): panic recovered; ", rr)
		}
	}()

	log.Debugf("HttpServerData::ServeHTTP(): begin")
	
	r.ParseForm()
	
    log.Debug("HttpServerData::ServeHTTP(): path   = ", r.URL.Path)
    log.Debug("HttpServerData::ServeHTTP(): method = ", r.Method)
    log.Debug("HttpServerData::ServeHTTP(): addr   = ", r.RemoteAddr)
	
    if r.Method == http.MethodGet {
	    f := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
    	
    	log.Debugf("HttpServerData::ServeHTTP(): f = %q", f)
		
		if len(f) < 3 {
			http.NotFound(w, r)
		} else {
			if server.isValidAPIKey(f[1]) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Info("HttpServerData::ServeHTTP(): err = ", err)
					http.NotFound(w, r)
				} else {
					log.Debugf("HttpServerData::ServeHTTP(): body = %+v", string(body[:]))
					var location Location
			
					err := json.Unmarshal(body, &location)
					if err != nil {
						log.Info("HttpServerData::ServeHTTP(): unmarshal err = ", err)
						http.NotFound(w, r)
					}
			
					location.dataId = f[2]

					log.Debugf("HttpServerData::ServeHTTP(): location = %+v", location)

					server.sendLocation(location)
				}
			} else {
				log.Infof("HttpServerData::ServeHTTP(): invalid API key '%s'", f[1])
				http.NotFound(w, r)
			}
		}
    } else {
    	log.Infof("HttpServerData::ServeHTTP(): not GET request")
	    http.NotFound(w, r)
    }

	log.Debugf("HttpServerData::ServeHTTP(): end")
}
//
//
func (server *HttpServerData) sendLocation(location Location)  {
    // send the location
    server.dispatcher.httpLocation <- location
}
//
//
func (server *HttpServerData) isValidAPIKey(apikey string) bool {
	for _, key := range server.config.apikeys {
		if key == apikey {
			return true
		}
	}

	return false
}