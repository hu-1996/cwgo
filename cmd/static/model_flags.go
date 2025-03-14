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
	"fmt"
	"strings"

	"github.com/hu-1996/cwgo/config"
	"github.com/hu-1996/cwgo/pkg/consts"
	"github.com/urfave/cli/v2"
)

func modelFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: consts.DSN, Usage: "Specify the database source name. (https://gorm.io/docs/connecting_to_the_database.html)", Value: "", DefaultText: "", Action: func(context *cli.Context, s string) error {
			if len(s) == 0 {
				return fmt.Errorf("dsn cannot be empty")
			}
			return nil
		}},
		&cli.StringFlag{Name: consts.DBType, Usage: "Specify database type. (mysql or sqlserver or sqlite or postgres)", Value: string(consts.MySQL), DefaultText: string(consts.MySQL), Action: func(context *cli.Context, s string) error {
			if _, ok := config.OpenTypeFuncMap[consts.DataBaseType(strings.ToLower(s))]; !ok {
				return fmt.Errorf("unknow db type %s (support mysql || postgres || sqlite || sqlserver for now)", s)
			}
			return nil
		}},
		&cli.StringFlag{Name: consts.OutDir, Usage: "Specify output directory", Value: consts.DefaultDbOutDir, DefaultText: consts.DefaultDbOutDir},
		&cli.StringFlag{Name: consts.OutFile, Usage: "Specify output filename", Value: consts.DefaultDbOutFile, DefaultText: consts.DefaultDbOutFile},
		&cli.StringSliceFlag{Name: consts.Tables, Usage: "Specify databases tables"},
		&cli.StringSliceFlag{Name: consts.ExcludeTables, Usage: "Specify exclude tables"},
		&cli.BoolFlag{Name: consts.UnitTest, Usage: "Specify generate unit test", Value: false, DefaultText: "false"},
		&cli.BoolFlag{Name: consts.OnlyModel, Usage: "Specify only generate model code", Value: false, DefaultText: "false"},
		&cli.StringFlag{Name: consts.ModelPkgName, Usage: "Specify model package name", Value: "", DefaultText: ""},
		&cli.BoolFlag{Name: consts.Nullable, Usage: "Specify generate with pointer when field is nullable", Value: false, DefaultText: "false"},
		&cli.BoolFlag{Name: consts.Signable, Usage: "Specify detect integer field's unsigned type, adjust generated data type", Value: false, DefaultText: "false"},
		&cli.BoolFlag{Name: consts.TypeTag, Usage: "Specify generate field with gorm column type tag", Value: false, DefaultText: "false"},
		&cli.BoolFlag{Name: consts.IndexTag, Usage: "Specify generate field with gorm index tag", Value: false, DefaultText: "false"},
		&cli.StringFlag{Name: consts.SQLDir, Usage: "Specify a sql file or directory", Value: "", DefaultText: ""},
	}
}
