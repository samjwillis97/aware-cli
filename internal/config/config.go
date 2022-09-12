package config

import (
	"fmt"
	"os"
	"path"
	"strings"

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
	Password     string // Going to be Stored in Config File ... Not ideal
	AuthProvider string
	Token        string // Gonig to be stored in Config File ... Not ideal
	Force        bool
}

type AwareCLIConfigGenerator struct {
	userCfg *AwareCLIConfig
	value   struct {
		server       string
		organisation string
		authProvider string
		login        string
		password     string
		token        string
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
    // TODO FIX THIS
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
	if err := c.configureServer(); err != nil {
		return "", err
	}

	if err := c.configureAuthProvider(); err != nil {
		return "", err
	}

	if err := c.configureLogin(); err != nil {
		return "", err
	}

	if err := c.configureOrganisation(); err != nil {
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
	config.Set("token", c.value.token)

    // TODO: Password? Maybe

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
				Help:    "This is the URL to the AWARE backend",
			},
			Validate: func(val interface{}) error {
				// TODO: Regex check URL
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

func (c *AwareCLIConfigGenerator) configureAuthProvider() error {
	c.value.authProvider = c.userCfg.AuthProvider
	c.value.token = c.userCfg.Token

	if c.userCfg.Token == "" {
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
				Help:    "This is the authentication provider you would like to login with.",
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
	}

	return nil
}

func (c *AwareCLIConfigGenerator) configureLogin() error {
	var qs []*survey.Question

	c.value.login = c.userCfg.Login
	c.value.password = c.userCfg.Password
	c.value.token = c.userCfg.Token

	if c.userCfg.Token == "" {
		if c.userCfg.Login == "" {
			qs = append(qs, &survey.Question{
				Name: "login",
				Prompt: &survey.Input{
					Message: "AWARE Login Email:",
					Help:    "This is your login username/email to the AWARE backend",
				},
				Validate: func(val interface{}) error {
					// TODO: Validate with email regex
					return nil
				},
			})
		}

		if c.userCfg.Password == "" {
			qs = append(qs, &survey.Question{
				Name: "password",
				Prompt: &survey.Password{
					Message: "AWARE Login Password:",
					Help:    "This is your login password to the AWARE backend",
				},
			})
		}
		if len(qs) > 0 {
			ans := struct {
				Login    string
				Password string
			}{}

			if err := survey.Ask(qs, &ans); err != nil {
				return err
			}

			c.value.login = ans.Login
			c.value.password = ans.Password
		}

		if ret, err := c.getLoginToken(c.value.server, c.value.login, c.value.password, c.value.authProvider); err != nil {
			return err
		} else {
			c.value.token = ret
			return nil
		}
	} else {
		// TODO: Validate Token
	}
	return nil
}

func (c *AwareCLIConfigGenerator) configureOrganisation() error {
    if c.userCfg.Organisation == "" {
        var organisationLabels []string
        var organisationIDs []string
        if ret, err := c.getOrganisations(c.value.server, c.value.token); err != nil {
            return err
        } else {
            for _, val := range ret {
                organisationIDs = append(organisationIDs, val.ID)
                organisationLabels = append(organisationLabels, val.Name)
            }
        }

        if len(organisationIDs) == 0 {
            // TODO: Error Out - No Orgs
            fmt.Println("NO ORGS")
            return nil
        }

        qs := &survey.Select{
            Message: "Default Organisation:",
            Help:    "This is the defualt organisation to use within AWARE.",
            Options: organisationLabels,
            Default: organisationLabels[0],
        }

        var selectedOrganisationLabel string
        if err := survey.AskOne(qs, &selectedOrganisationLabel); err != nil {
            return err
        }

        for i, val := range organisationLabels {
            if selectedOrganisationLabel == val {
                c.value.organisation = organisationIDs[i]
                break
            }
        }

        if c.value.organisation == "" {
            // TODO Error
            fmt.Println("ERROR NO ORG")
            return nil
        }
    } else {
        // TODO: Validate Org exists
    }
	return nil
}

func (c *AwareCLIConfigGenerator) getAuthProviders(server string) ([]*aware.AuthProvider, error) {
	s := utils.ShowLoading("Getting Authentication Providers...")
	defer s.Stop()

	server = strings.TrimRight(server, "/")

	c.awareClient = aware.NewClient(aware.Config{
		Server:   server,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	if ret, err := c.awareClient.GetAllAuthProviders(); err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

func (c *AwareCLIConfigGenerator) getLoginToken(server, login, password, authType string) (string, error) {
	s := utils.ShowLoading("Verifying Login Details and Getting JWT...")
	defer s.Stop()

	server = strings.TrimRight(server, "/")

	c.awareClient = aware.NewClient(aware.Config{
		Server:   server,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	if ret, err := c.awareClient.Login(login, password, authType); err != nil {
		return "", err
	} else {
		return ret.AccessToken, nil
	}
}

func (c *AwareCLIConfigGenerator) getOrganisations(server, token string) ([]*aware.Organisation, error) {
	s := utils.ShowLoading("Getting Organisations...")
	defer s.Stop()

	server = strings.TrimRight(server, "/")

	c.awareClient = aware.NewClient(aware.Config{
		Server:   server,
        Token: token,
		Insecure: true,
		Debug:    viper.GetBool("debug"),
	})

	if ret, err := c.awareClient.GetAllOrganisations(); err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}
