package command

import (
	"bytes"
	"fmt"

	"github.com/spf13/pflag"
)

type CmdData interface {
	Name() string
	ShortDescr() string
	Run( args []string)
	usageStr() string
}

func MapNames[C CmdData](cmds []C) map[string]C {
	cmdMap := make(map[string]C, len(cmds))
	for _, cmd := range cmds {
		cmdMap[cmd.Name()] = cmd
	}
	return cmdMap
}

type BaseArgs struct {
	ConfigPath string
}

func InitBaseFlags(argData *BaseArgs) *pflag.FlagSet {
	flags := pflag.NewFlagSet("base flagset", pflag.PanicOnError)
	flags.StringVar(
		&argData.ConfigPath,
		"config",
		"",
		"The path to a Prcrastinate configuration file",
	)
	return flags
}

// func (args *BaseArgs) applyToConfig(cfg *platform.Config) {
// 	if len(args.LogPath) > 0 {
// 		cfg.LogPath = args.LogPath
// 	}
// 	if args.Verbose {
// 		cfg.Verbosity = platform.Verbose
// 	}
// }

func FmtHelp(c CmdData) string {
	var outputBuf bytes.Buffer

	fmt.Fprintf(&outputBuf, "prcr %s\n", c.Name())
	fmt.Fprintf(&outputBuf, "%s\n\n", c.ShortDescr())

	return outputBuf.String()
}
