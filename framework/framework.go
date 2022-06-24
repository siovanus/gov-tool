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

package framework

import (
	"github.com/polynetwork/gov-tool/log"
	"github.com/polynetwork/gov-tool/zion"
)

//Default TestFramework instance
var TFramework = NewFramework()

type Case func(ctx *FrameworkContext) bool

//Framework manage case and run case
type Framework struct {
	casesMap map[string]Case
}

func NewFramework() *Framework {
	return &Framework{
		casesMap: make(map[string]Case, 0),
	}
}

//RegCase register a case to framework
func (this *Framework) RegCase(name string, c Case) {
	this.casesMap[name] = c
}

func (this *Framework) RunCase(ctx *FrameworkContext, name string) {
	c, found := this.casesMap[name]
	if !found {
		panic("can not find this case name")
	}
	ok := c(ctx)
	if !ok {
		log.Error("case failed")
	} else {
		log.Info("case success")
	}
}

//FrameworkContext is the context for case
type FrameworkContext struct {
	Z *zion.ZionTools
}

func NewFrameworkContext(z *zion.ZionTools) *FrameworkContext {
	return &FrameworkContext{
		Z: z,
	}
}
