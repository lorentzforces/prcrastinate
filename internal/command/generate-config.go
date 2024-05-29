package command

import (
	"fmt"
	"os"
	"path/filepath"
	"prcrastinate/internal/platform"

	"github.com/BurntSushi/toml"
)

type GenerateConfig struct{}

func (cmd GenerateConfig) Name() string {
	return "generate-config"
}

func (cmd GenerateConfig) ShortDescr() string {
	return "Generate a Prcrastinate config file"
}

func (cmd GenerateConfig) usageStr() string {
	return "PLACEHOLDER usage for [generate-config]"
}

type ConfigGenerateArgs struct {
	BaseArgs
}

func (cmd GenerateConfig) Run(args []string) {
	parsedArgs := new(ConfigGenerateArgs)
	baseFlags := InitBaseFlags(&parsedArgs.BaseArgs)
	baseFlags.Parse(args)

	// Error usually means that we didn't find the file. If something will prevent us writing the
	// file, we'll encounter that later anyway.
	configPath, foundFile, _  := platform.GetConfigPath(parsedArgs.ConfigPath)
	if foundFile {
		platform.FailOut(fmt.Sprintf("Resolved config file path already exists: %s", configPath))
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o_777); err != nil {
		platform.FailOut(fmt.Sprintf(
			"Could not create directory structure for config file: %s\n%s",
			configPath,
			err.Error(),
		))
	}

	configFile, err := os.Create(configPath)
	defer configFile.Close()
	if err != nil {
		platform.FailOut(fmt.Sprintf(
			"Could not open file for writing: %s\n%s",
			configPath,
			err.Error(),
		))
	}

	_, err = configFile.WriteString("# this file was generated by Prcrastinate\n")
	if err != nil {
		platform.FailOut(err.Error())
	}

	tomlEncoder := toml.NewEncoder(configFile)
	err = tomlEncoder.Encode(platform.DefaultConfig())
	if err != nil {
		platform.FailOut(fmt.Sprintf(
			"Failed to write default config to %s, but the file was created successfully. You " +
				"may wish to delete this file.\n%s",
			configPath,
			err.Error(),
		))
	}
}
