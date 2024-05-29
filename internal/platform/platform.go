package platform

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

const configEnvVar = "PRCR_CONFIG"
const defaultConfigPath =  "~/.config/prcrastinate/config.toml"

func FailOut(msg string) {
	fmt.Fprintln(os.Stderr, "ERROR: " + msg)
	os.Exit(1)
}

func Assert(condition bool, more any) {
	if condition { return }
	panic(fmt.Sprintf("Assertion violated: %s", more))
}

type Config struct {
	Token string
	Endpoint string
}

func DefaultConfig() Config {
	return Config{
		Token: "",
		Endpoint: "https://api.github.com/graphql",
	}
}

// Reads config using the provided path argument.
// This is the intended method for commands to obtain configuration information: any error
// encountered while reading configuration will result in an error exit.
func ReadConfigFromPath(path string) Config {
	path, foundFile, err := GetConfigPath(path)
	if err != nil {
		FailOut(err.Error())
	}

	if !foundFile {
		return DefaultConfig()
	}

	f, err := os.Open(path)
	if err != nil {
		FailOut(err.Error())
	}

	return ReadConfig(f)
}

func ReadConfig(r io.Reader) Config {
	configData := DefaultConfig()
	decoder := toml.NewDecoder(r)

	_, err := decoder.Decode(&configData)
	if err != nil {
		FailOut(err.Error())
	}
	return configData
}

// Resolves the config from the specified path, if any. If the passed path is empty, will use:
//   - value of the PRCR_CONFIG environment variable
//   - the default config path on the platform
// For now, prcrastinate uses the default XDG_CONFIG directory: ~/.config/prcrastinate/config.toml
func GetConfigPath(configPath string) (path string, foundFile bool, err error) {
	fileSpecified := true
	envPath := os.Getenv(configEnvVar)

	if len(configPath) == 0 && len(envPath) > 0 {
		configPath = envPath
	}

	if len(configPath) == 0 {
		fileSpecified = false
		configPath = replaceTilde(defaultConfigPath)
	}

	configFile, err := os.Open(configPath)
	defer func() {
		configFile.Close()
	}()

	// no config file exists, so just use the default
	if err != nil && !fileSpecified {
		// TODO: perhaps log that default configuration was used?
		return configPath, false, nil
	}
	// config file was specified, but we couldn't read it
	if err != nil {
		return configPath, false, err
	}
	 // in theory if above succeeded this can't fail
	configStat, err := configFile.Stat()
	Assert(err == nil, err)

	if configStat.IsDir() {
		return configPath,
			true,
			fmt.Errorf("Resolved config file is a directory: \"%s\"", configPath)
	}

	return configPath, true, nil
}

func replaceTilde(s string) string {
	if strings.HasPrefix(s, "~") {
		usr, err := user.Current()
		if err != nil {
			FailOut(err.Error())
		}
		s = strings.Replace(s, "~", usr.HomeDir, 1)
	}
	return s
}

func ValidateConfig(config *Config) error {
	var errs []error
	if emptyToken, _ := regexp.MatchString(`^\s*$`, config.Token); emptyToken {
		errs = append(errs, fmt.Errorf("No valid token value provided"))
	}

	// TODO: validate the rest

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
