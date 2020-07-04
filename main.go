/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	_ "expvar"
	"infini.sh/framework"
	"infini.sh/framework/core/module"
	"infini.sh/framework/core/orm"
	"infini.sh/framework/core/util"
	"infini.sh/framework/modules"
	"infini.sh/proxy/config"
	"infini.sh/proxy/model"
	"infini.sh/proxy/plugin"
	"infini.sh/proxy/ha"
)

func main() {
	terminalHeader := ("___  ____ ____ _  _ _   _\n")
	terminalHeader += ("|__] |__/ |  |  \\/   \\_/\n")
	terminalHeader += ("|    |  \\ |__| _/\\_   |\n")

	terminalFooter := ("                         |    |                \n")
	terminalFooter += ("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) \n")
	terminalFooter += (" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| \n")
	terminalFooter += (" ____/                             ___/        \n")

	app := framework.NewApp("proxy", "An elasticsearch proxy written in golang.",
		util.TrimSpaces(config.Version), util.TrimSpaces(config.LastCommitLog), util.TrimSpaces(config.BuildDate), terminalHeader, terminalFooter)

	app.Init(nil)

	defer app.Shutdown()

	app.Start(func() {
        ha.SetupVIP()

		//load core modules first
		modules.Register()

		module.RegisterUserPlugin(plugin.ProxyPlugin{})

		//start each module, with enabled provider
		module.Start()

	}, func() {
		orm.RegisterSchema(&model.Request{})

	})

}
