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

package config

import (
	"github.com/hu-1996/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

type ApiArgument struct {
	ProjectPath  string
	HertzRepoUrl string
}

func NewApiArgument() *ApiArgument {
	return &ApiArgument{}
}

func (c *ApiArgument) ParseCli(ctx *cli.Context) error {
	c.ProjectPath = ctx.String(consts.ProjectPath)
	c.HertzRepoUrl = ctx.String(consts.HertzRepoUrl)
	return nil
}
