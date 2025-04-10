/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"os"

	"github.com/hu-1996/cwgo/hertz"
	"github.com/hu-1996/cwgo/pkg/curd/doc/mongo/plugin"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/protoc"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"
	"github.com/hu-1996/cwgo/cmd/static"
	"github.com/hu-1996/cwgo/tpl"
)

func main() {
	tpl.RegisterTemplateFunc()

	// run cwgo as hz plugin mode
	hertz.PluginMode()
	// run cwgo as kitex plugin mode
	kitexPluginMode()
	// run cwgo as mongo plugin mode
	plugin.MongoPluginMode()

	tpl.Init()
	cli := static.Init()

	err := cli.Run(os.Args)
	if err != nil {
		logs.Errorf("%v\n", err)
	}
}

func kitexPluginMode() {
	mode := os.Getenv(kargs.EnvPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		// run as a plugin
		switch mode {
		case thriftgo.PluginName:
			os.Exit(thriftgo.Run())
		case protoc.PluginName:
			os.Exit(protoc.Run())
		}
	}
}
