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

package client

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/internal_pkg/generator"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/hu-1996/cwgo/config"
	"github.com/hu-1996/cwgo/pkg/common/utils"
	"github.com/hu-1996/cwgo/pkg/consts"
	"github.com/hu-1996/cwgo/tpl"
)

func convertKitexArgs(sa *config.ClientArgument, kitexArgument *kargs.Arguments) (err error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	kitexArgument.ModuleName = sa.GoMod
	kitexArgument.ServiceName = sa.ServerName
	kitexArgument.Includes = sa.SliceParam.ProtoSearchPath
	kitexArgument.Version = kitex.Version
	kitexArgument.RecordCmd = os.Args
	kitexArgument.ThriftOptions = append(kitexArgument.ThriftOptions,
		"naming_style=golint",
		"ignore_initialisms",
		"gen_setter",
		"gen_deep_equal",
		"compatible_names",
		"frugal_tag",
	)
	kitexArgument.IDL = sa.IdlPath

	f.BoolVar(&kitexArgument.NoFastAPI, "no-fast-api", false, "Generate codes without injecting fast method.")
	f.StringVar(&kitexArgument.Use, "use", "",
		"Specify the kitex_gen package to import when generate server side codes.")
	f.BoolVar(&kitexArgument.GenerateInvoker, "invoker", false,
		"Generate invoker side codes when service name is specified.")
	f.StringVar(&kitexArgument.IDLType, "type", "unknown", "Specify the type of IDL: 'thrift' or 'protobuf'.")
	f.Var(&kitexArgument.ThriftOptions, "thrift", "Specify arguments for the thrift go compiler.")
	f.DurationVar(&kitexArgument.ThriftPluginTimeLimit, "thrift-plugin-time-limit", generator.DefaultThriftPluginTimeLimit, "Specify thrift plugin execution time limit.")
	f.Var(&kitexArgument.ThriftPlugins, "thrift-plugin", "Specify thrift plugin arguments for the thrift compiler.")
	f.Var(&kitexArgument.ProtobufOptions, "protobuf", "Specify arguments for the protobuf compiler.")
	f.BoolVar(&kitexArgument.CombineService, "combine-service", false,
		"Combine services in root thrift file.")
	f.BoolVar(&kitexArgument.CopyIDL, "copy-idl", false,
		"Copy each IDL file to the output path.")
	f.StringVar(&kitexArgument.ExtensionFile, "template-extension", kitexArgument.ExtensionFile,
		"Specify a file for template extension.")
	f.BoolVar(&kitexArgument.FrugalPretouch, "frugal-pretouch", false,
		"Use frugal to compile arguments and results when new clients and servers.")
	f.BoolVar(&kitexArgument.Record, "record", false, "Record Kitex cmd into kitex-all.sh.")
	f.StringVar(&kitexArgument.GenPath, "gen-path", generator.KitexGenPath,
		"Specify a code gen path.")
	f.Var(&kitexArgument.ProtobufPlugins, "protobuf-plugin", "Specify protobuf plugin arguments for the protobuf compiler.(plugin_name:options:out_dir)")
	f.StringVar(&kitexArgument.Protocol, "protocol", "", "Specify a protocol for codec.")
	f.Var(&kitexArgument.Hessian2Options, "hessian2", "Specify arguments for the hessian2 codec.")

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, `Version %s
Usage: %s [flags] IDL

Flags:
`, kitexArgument.Version, os.Args[0])
		f.PrintDefaults()
		os.Exit(1)
	}

	err = f.Parse(utils.StringSliceSpilt(sa.SliceParam.Pass))
	if err != nil {
		return
	}

	kitexArgument.GenerateMain = false

	// Non-standard template
	if strings.HasSuffix(sa.Template, consts.SuffixGit) {
		err = utils.GitClone(sa.Template, path.Join(tpl.KitexDir, consts.Client))
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(sa.Template)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.KitexDir, consts.Client, gitPath)
		if err = utils.GitCheckout(sa.Branch, gitPath); err != nil {
			return err
		}
		kitexArgument.TemplateDir = gitPath
	} else {
		if len(sa.Template) != 0 {
			kitexArgument.TemplateDir = sa.Template
		} else {
			kitexArgument.TemplateDir = path.Join(tpl.KitexDir, consts.Client, consts.Standard)
		}
	}

	return checkKitexArgs(kitexArgument)
}

func checkKitexArgs(a *kargs.Arguments) (err error) {
	// check IDL
	a.IDLType, err = utils.GetIdlType(a.IDL, consts.Protobuf)
	if err != nil {
		return err
	}

	// check service name
	if a.ServiceName == "" {
		if a.Use != "" {
			log.Warn("-use must be used with -service")
			os.Exit(2)
		}
	}

	gopath, err := utils.GetGOPATH()
	if err != nil {
		return fmt.Errorf("get gopath failed: %s", err)
	}
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	gosrc := filepath.Join(gopath, "src")
	gosrc, err = filepath.Abs(gosrc)
	if err != nil {
		log.Warn("Get GOPATH/src path failed:", err.Error())
		os.Exit(1)
	}
	curpath, err := filepath.Abs(".")
	if err != nil {
		log.Warn("Get current path failed:", err.Error())
		os.Exit(1)
	}

	if strings.HasPrefix(curpath, gosrc) {
		if a.PackagePrefix, err = filepath.Rel(gosrc, curpath); err != nil {
			log.Warn("Get GOPATH/src relpath failed:", err.Error())
			os.Exit(1)
		}
		a.PackagePrefix = filepath.Join(a.PackagePrefix, generator.KitexGenPath)
	} else {
		if a.ModuleName == "" {
			log.Warn("Outside of $GOPATH. Please specify a module name with the '-module' flag.")
			os.Exit(1)
		}
	}

	if a.ModuleName != "" {
		module, path, ok := utils.SearchGoMod(curpath, true)
		if ok {
			// go.mod exists
			if module != a.ModuleName {
				log.Warnf("The module name given by the '-module' option ('%s') is not consist with the name defined in go.mod ('%s' from %s)\n",
					a.ModuleName, module, path)
				os.Exit(1)
			}
			if a.PackagePrefix, err = filepath.Rel(path, curpath); err != nil {
				log.Warn("Get package prefix failed:", err.Error())
				os.Exit(1)
			}
			a.PackagePrefix = filepath.Join(a.ModuleName, a.PackagePrefix, generator.KitexGenPath)
		} else {
			if err = utils.InitGoMod(a.ModuleName); err != nil {
				log.Warn("Init go mod failed:", err.Error())
				os.Exit(1)
			}
			a.PackagePrefix = filepath.Join(a.ModuleName, generator.KitexGenPath)
		}
	}

	if a.Use != "" {
		a.PackagePrefix = a.Use
	}
	a.OutputPath = curpath
	a.PackagePrefix = strings.ReplaceAll(a.PackagePrefix, consts.BackSlash, consts.Slash)
	return nil
}
