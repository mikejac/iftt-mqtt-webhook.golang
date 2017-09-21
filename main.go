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
	 "time"
	 "os"
	 "github.com/docopt/docopt-go"
	 "github.com/kardianos/service"
	 "github.com/mikejac/log.golang"
 )
 
 type Program struct {
	 exit chan	bool
 }
 
 var (
	 Config	*DispatcherConfiguration
	 usage string = `IFTTT-MQTT Webhook.
 
 Usage:
   iftt-mqtt-webhook <configfile> [--install] [--debug]
   iftt-mqtt-webhook (--start|--stop|--restart|--uninstall)
   iftt-mqtt-webhook -h | --help
   iftt-mqtt-webhook --version
 
 Options:
   -h --help     Show this screen.
   --version     Show version.`
 )
 
 //
 //
 func main() {
	 arguments, _ := docopt.Parse(usage, nil, true, "IFTTT-MQTT Webhook 1.0", false)
 
	 if arguments["--debug"].(bool) {
		 log.EnableDebugLog(true)
	 }
 
	 /******************************************************************************************************************
	  * prepare our service stuff
	  *
	  */
	 svcConfig := &service.Config{
		 Name:        "iftt-mqtt-webhook",
		 DisplayName: "iftt-mqtt-webhook",
		 Description: "IFTTT-MQTT Webhook ver. 1.0",
	 }
 
	 Config = NewConfig()
	 
	 if arguments["<configfile>"] != nil {
		 if err := Config.ReadConfig(arguments["<configfile>"].(string)); err != nil {
			 log.Infof("error: %s", err.Error())
			 log.Flush()
			 return
		 }
		 
		 svcConfig.Arguments = append(svcConfig.Arguments, arguments["<configfile>"].(string))
	 } else {
		 log.Infof("error: no configuration file specified on command line")
		 log.Flush()
		 return
	 }
 
	 prg := &Program{}
	 s, err := service.New(prg, svcConfig)
	 if err != nil {
		 log.Info("error: ", err)
		 log.Flush()
		 return
	 }
 
	 if arguments["--start"].(bool) {
		 err := service.Control(s, "start")
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 } else if arguments["--stop"].(bool) {
		 err := service.Control(s, "stop")
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 } else if arguments["--restart"].(bool) {
		 err := service.Control(s, "restart")
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 } else if arguments["--install"].(bool) {
		 err := service.Control(s, "install")
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 } else if arguments["--uninstall"].(bool) {
		 err := service.Control(s, "uninstall")
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 } else {
		 /******************************************************************************************************************
		  * now we're ready to run. the action now takes place in 'run()'
		  *
		  */
		 err = s.Run()
		 if err != nil {
			 log.Info("error: ", err)
		 }
	 }
 
	 log.Flush()
 }
 
 //
 //
 func (p *Program) run() {
	 log.Debug("Program::run(): begin")
	 
	 disp := NewDispatcher(Config, p.exit)
	 if disp == nil {
		 log.Infof("failed to allocate new dispatcher")
		 log.Flush()
		 
		 os.Exit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	 }
 
	 /******************************************************************************************************************
	  * main loop 
	  *
	  */
	 if err := disp.Run(); err != nil {
		 log.Info("error: ", err)
		 log.Flush()
 
		 os.Exit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	 }
	 
	 log.Debug("Program::run(): end")
	 log.Flush()
 }
 
 //
 //
 func (p *Program) Start(s service.Service) error {
	 // start should not block. Do the actual work async.
	 log.Debug("Program::Start(): begin")
 
	 if service.Interactive() {
		 log.Debug("Program::Start(): running in terminal")
	 } else {
		 log.Debug("Program::Start(): running under service manager")
	 }
 
	 p.exit = make(chan bool) // our exit-signal
 
	 // let's get going
	 go p.run()
 
	 log.Debug("Program::Start(): end")
	 
	 return nil
 }
 
 //
 //
 func (p *Program) Stop(s service.Service) error {
	 // stop should not block. Return within a few seconds
	 log.Debug("Program::Stop(): begin")
 
	 p.exit <- true
 
	 // the run() function needs a bit of time to catch up
	 time.Sleep(1 * time.Second)
 
	 log.Debug("Program::Stop(): end")
	 log.Flush()
	 
	 return nil
 }
