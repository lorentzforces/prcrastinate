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
const DefaultConfigPath =  "~/.config/prcrastinate/config.toml"

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

func GetConfigEnvVarPath() string {
	// TODO: consider logging a warning if the env var is set but blank
	envVarPath := os.Getenv(configEnvVar)
	if isBlank(envVarPath) {
		return ""
	}
	return envVarPath
}

// Reads config using the provided path argument.
// This is the intended method for commands to obtain configuration information: any error
// encountered while reading configuration will result in an error exit.
// Will read the first non-blank path specified out of the following list:
//   - value of the PRCR_CONFIG environment variable
//   - the path provided as an argument to this function
//   - the default config path on the platform
// For now, prcrastinate uses the default XDG_CONFIG directory: ~/.config/prcrastinate/config.toml
func ReadConfigFromPath(path string) Config {
	envVarPath := GetConfigEnvVarPath()
	if len(envVarPath) > 0 {
		path = envVarPath
	}

	path, foundFile, err := GetDefaultablePath(path, DefaultConfigPath)
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

func GetDefaultablePath(givenPath string, defaultPath string) (path string, foundFile bool, err error) {
	givenPath = replaceTilde(givenPath)
	defaultPath = replaceTilde(defaultPath)

	workingPath := givenPath

	fileSpecified := true
	if isBlank(givenPath) {
		fileSpecified = false
		workingPath = defaultPath
	}

	file, err := os.Open(workingPath)
	defer file.Close()

	// no file exists, so just use the default
	if err != nil && !fileSpecified {
		return workingPath, false, nil
	}
	// file was specified, but we couldn't read it
	if err != nil {
		return workingPath, false, err
	}
	// in theory if above succeeded this can't fail
	fileStat, err := file.Stat()
	Assert(err == nil, err)

	if fileStat.IsDir() {
		return workingPath,
			true,
			fmt.Errorf("Resolved file path is a directory: \"%s\"", workingPath)
	}

	return workingPath, true, nil
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
	if isBlank(config.Token) {
		errs = append(errs, fmt.Errorf("No valid token value provided"))
	}

	// TODO: validate the rest

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func isBlank(str string) bool {
	isBlank, _ := regexp.MatchString(`^\s*$`, str)
	return isBlank
}
