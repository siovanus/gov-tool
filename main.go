/*
* Copyright (C) 2020 The poly network Authors
* This file is part of The poly network library.
*
* The poly network is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The poly network is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
* You should have received a copy of the GNU Lesser General Public License
* along with The poly network . If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	_ "github.com/polynetwork/gov-tool/cases"
	"github.com/polynetwork/gov-tool/config"
	"github.com/polynetwork/gov-tool/framework"
	"github.com/polynetwork/gov-tool/log"
	"github.com/polynetwork/gov-tool/zion"
	"time"
)

var (
	Config   string
	CaseName string
)

func init() {
	flag.StringVar(&Config, "cfg", "./config.json", "Config of gov-tool")
	flag.StringVar(&CaseName, "c", "", "Case name to run")
	flag.Parse()
}

func main() {
	log.InitLog(2) //init log module
	defer time.Sleep(time.Second)

	err := config.DefConfig.Init(Config)
	if err != nil {
		log.Error("DefConfig.Init error:%s", err)
		return
	}

	z := zion.NewZionTools(config.DefConfig.JsonRpcAddress)
	ctx := framework.NewFrameworkContext(z)
	framework.TFramework.RunCase(ctx, CaseName)
}
