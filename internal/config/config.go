package config

import (
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

    "ampaware.com/cli/internal/utils"
)

const (
	ConfigFileName = "aware_config"
	ConfigFileType = "yaml"
)

type AwareCLIConfig struct {
	Server       string
	Organisation string
	Login        string
	AuthProvider string
    Force bool
}

func Exists(file string) bool {
	if file == "" {
		return false
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetConfigDirectory() (string, error) {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home != "" {
		return path.Join(home, ".config", "aware"), nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".config", "aware"), nil
}

func CheckForToken() {
	if os.Getenv("AWARE_JWT") != "" {
		return
	}

	// TODO: Fixup, -> Inform
	fmt.Printf("JWT Required.")

	os.Exit(1)
}

func (c *AwareCLIConfig) Generate() (string, error) {
    cfgDir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}

	if err := func() error {
		s := utils.ShowLoading("Creating new configuration...")
		defer s.Stop()

		return createFile(cfgDir, fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileType))
	}(); err != nil {
		return "", err
	}

	return c.writeToFile(cfgDir)
}

func createFile(path, name string) error {
	const perm = 0o700

	if !Exists(path) {
		if err := os.MkdirAll(path, perm); err != nil {
			return err
		}
	}

	file := fmt.Sprintf("%s/%s", path, name)
	if Exists(file) {
		if err := os.Rename(file, file+".bkp"); err != nil {
			return err
		}
	}
	_, err := os.Create(file)

	return err
}

func (c *AwareCLIConfig) writeToFile(path string) (string, error) {
    config := viper.New()

    config.AddConfigPath(path)
    config.SetConfigName(ConfigFileName)
    config.SetConfigType(ConfigFileType)

    config.Set("server", c.Server)
    config.Set("organisation", c.Organisation)
    config.Set("login", c.Login)
    config.Set("authProvider", c.AuthProvider)

    if err := config.WriteConfig(); err != nil {
		return "", err
	}

    return fmt.Sprintf("%s%s%s.%s", path, string(os.PathSeparator), ConfigFileName, ConfigFileType), nil
}
