package hertz

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/hertz/cmd/hz/meta"
	"github.com/cloudwego/hertz/cmd/hz/thrift"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/hu-1996/cwgo/config"
	"github.com/hu-1996/cwgo/hertz/generator"
	"github.com/hu-1996/cwgo/hertz/protobuf"
)

func GenerateLayout(args *config.HzArgument) error {
	lg := &generator.LayoutGenerator{
		TemplateGenerator: generator.TemplateGenerator{
			OutputDir: args.OutDir,
			Excludes:  args.Excludes,
		},
	}

	layout := generator.Layout{
		GoModule:        args.Gomod,
		ServiceName:     args.ServiceName,
		UseApacheThrift: args.IdlType == meta.IdlThrift,
		HasIdl:          0 != len(args.IdlPaths),
		ModelDir:        args.ModelDir,
		HandlerDir:      args.HandlerDir,
		RouterDir:       args.RouterDir,
		NeedGoMod:       args.NeedGoMod,
	}

	if args.CustomizeLayout == "" {
		// generate by default
		err := lg.GenerateByService(layout)
		if err != nil {
			return fmt.Errorf("generating layout failed: %v", err)
		}
	} else {
		// generate by customized layout
		configPath, dataPath := args.CustomizeLayout, args.CustomizeLayoutData
		logs.Infof("get customized layout info, layout_config_path: %s, template_data_path: %s", configPath, dataPath)
		exist, err := util.PathExist(configPath)
		if err != nil {
			return fmt.Errorf("check customized layout config file exist failed: %v", err)
		}
		if !exist {
			return errors.New("layout_config_path doesn't exist")
		}
		lg.ConfigPath = configPath
		// generate by service info
		if dataPath == "" {
			err := lg.GenerateByService(layout)
			if err != nil {
				return fmt.Errorf("generating layout failed: %v", err)
			}
		} else {
			// generate by customized data
			err := lg.GenerateByConfig(dataPath)
			if err != nil {
				return fmt.Errorf("generating layout failed: %v", err)
			}
		}
	}

	err := lg.Persist()
	if err != nil {
		return fmt.Errorf("generating layout failed: %v", err)
	}
	return nil
}

func TriggerPlugin(args *config.HzArgument) error {
	if len(args.IdlPaths) == 0 {
		return nil
	}
	cmd, err := config.BuildPluginCmd(args)
	if err != nil {
		return fmt.Errorf("build plugin command failed: %v", err)
	}

	compiler, err := config.IdlTypeToCompiler(args.IdlType)
	if err != nil {
		return fmt.Errorf("get compiler failed: %v", err)
	}

	logs.Debugf("begin to trigger plugin, compiler: %s, idl_paths: %v", compiler, args.IdlPaths)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		out := strings.TrimSpace(string(buf))
		if !strings.HasSuffix(out, meta.TheUseOptionMessage) {
			return fmt.Errorf("plugin %s_gen_hertz returns error: %v, cause:\n%v", compiler, err, string(buf))
		}
	}

	// If len(buf) != 0, the plugin returned the log.
	if len(buf) != 0 {
		fmt.Println(string(buf))
	}
	logs.Debugf("end run plugin %s_gen_hertz", compiler)
	return nil
}

func PluginMode() {
	mode := os.Getenv(meta.EnvPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		switch mode {
		case meta.ThriftPluginName:
			plugin := new(thrift.Plugin)
			os.Exit(plugin.Run())
		case meta.ProtocPluginName:
			plugin := new(protobuf.Plugin)
			os.Exit(plugin.Run())
		}
	}
}
