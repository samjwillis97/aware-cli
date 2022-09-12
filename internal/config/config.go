package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"ampaware.com/cli/internal/utils"
	"ampaware.com/cli/pkg/aware"
)

const (
	ConfigFileName = "aware_config"
	ConfigFileType = "yaml"
)

type AwareCLIConfig struct {
	Server       string
	Organisation string
	Login        string
    Password string
	AuthProvider string
    Force bool
}

type AwareCLIConfigGenerator struct {
    userCfg *AwareCLIConfig
    value struct {
        server string
        organisation string
        authProvider string
        login string
        password string
    }
    awareClient *aware.Client
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

func NewAwareCLIConfigGenerator(cfg *AwareCLIConfig) *AwareCLIConfigGenerator {
    gen := AwareCLIConfigGenerator{
        userCfg: cfg,
    }

    return &gen
}

func (c *AwareCLIConfigGenerator) Generate() (string, error) {
    cfgDir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}

    // TODO: Handle Force?
    // TODO: Questionairre
    // Server
    // Select AuthProvider
    // Login
    // Select Org
    if err := c.configureServer(); err != nil {
        return "", err
    }

    if err := c.configureAuthProvider(); err != nil {
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

func (c *AwareCLIConfigGenerator) writeToFile(path string) (string, error) {
    config := viper.New()

    config.AddConfigPath(path)
    config.SetConfigName(ConfigFileName)
    config.SetConfigType(ConfigFileType)

    config.Set("server", c.value.server)
    config.Set("authProvider", c.value.authProvider)
    config.Set("login", c.value.login)
    config.Set("organisation", c.value.organisation)

    if err := config.WriteConfig(); err != nil {
		return "", err
	}

    return fmt.Sprintf("%s%s%s.%s", path, string(os.PathSeparator), ConfigFileName, ConfigFileType), nil
}

func (c *AwareCLIConfigGenerator) configureServer() error {
    var qs []*survey.Question

    c.value.server = c.userCfg.Server

    if c.userCfg.Server == "" {
        qs = append(qs, &survey.Question{
            Name: "server",
            Prompt: &survey.Input{
                Message: "AWARE API URL:",
                Help: "This is the URL to the AWARE backend",
            },
            Validate: func(val interface{}) error {
                // TODO: Call Auth Providers? check for response
                return nil
            },
        })

        var server string
        if err := survey.Ask(qs, &server); err != nil {
            return err
        }

        c.value.server = server
    } 

    return nil
}

func (c *AwareCLIConfigGenerator) getAuthProviders(server string) ([]*aware.AuthProvider, error) {
    s := utils.ShowLoading("Getting Authentication Providers...") 
    defer s.Stop()

    server = strings.TrimRight(server, "/")

    c.awareClient = aware.NewClient(aware.Config{
        Server: server,
        Insecure: true,
        Debug: viper.GetBool("debug"),
    })

    if ret, err := c.awareClient.GetAllAuthProviders(); err != nil {
        return nil, err
    } else {
        return ret, nil
    }
}

func (c *AwareCLIConfigGenerator) configureAuthProvider() error {
    c.value.authProvider = c.userCfg.AuthProvider

    var providerTypes []string
    var providerLabels []string

    if ret, err := c.getAuthProviders(c.value.server); err != nil {
        return err
    } else {
        for _, val := range ret {
            providerTypes = append(providerTypes, val.AuthType)
            providerLabels = append(providerLabels, val.Label)
        }
    }

    if c.userCfg.AuthProvider == "" {
        qs := &survey.Select{
            Message: "Authentication Provider:",
            Help: "This is the authentication provider you would like to login with.",
            Options: providerLabels,
            Default: providerLabels[0],
        }

        var authProviderLabel string
        if err := survey.AskOne(qs, &authProviderLabel); err != nil {
            return err
        }

        for i, val := range providerLabels {
            if authProviderLabel == val {
                c.value.authProvider = providerTypes[i]
                break
            }
        }

        if c.value.authProvider == "" {
            // TODO: Error
            fmt.Println("ERROR NO AUTH")
        }
    } else {
        // TODO: Validate Chosen Auth
    }

    return nil
}
