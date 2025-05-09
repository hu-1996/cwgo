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

package static

import (
	"github.com/hu-1996/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func apiFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  consts.ProjectPath,
			Usage: "Specify the project path.",
		},
		&cli.StringFlag{
			Name:        consts.HertzRepoUrl,
			Aliases:     []string{"r"},
			DefaultText: consts.HertzRepoDefaultUrl,
			Usage:       "Specify the url of the hertz repository you want",
		},
	}
}
