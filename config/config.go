package config

import (
	"fmt"
	"os"

	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

const (
	defaultLogLevel = "DEBUG"
	defaultLookback = "24h"
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	} `yaml:"log"`

	Notion struct {
		ApiToken       string `yaml:"api_token" env:"NOT_API_TOKEN" valid:"minstringlength(3)"`
		OrganisationID string `yaml:"organisation_id" env:"NOT_ORG_ID" valid:"minstringlength(3)"`
		Lookback       string `yaml:"lookback" env:"NOT_LOOKBACK"`
	} `yaml:"notion"`

	Microsoft struct {
		AppID          string `yaml:"app_id" env:"MS_APP_ID" valid:"minstringlength(3)"`
		SecretKey      string `yaml:"secret_key" env:"MS_SECRET_KEY" valid:"minstringlength(3)"`
		TenantID       string `yaml:"tenant_id" env:"MS_TENANT_ID" valid:"minstringlength(3)"`
		SubscriptionID string `yaml:"subscription_id" env:"MS_SUB_ID" valid:"minstringlength(3)"`

		DataCollection struct {
			Endpoint   string `yaml:"endpoint" env:"MS_DCR_ENDPOINT" valid:"minstringlength(3)"`
			RuleID     string `yaml:"rule_id" env:"MS_DCR_RULE" valid:"minstringlength(3)"`
			StreamName string `yaml:"stream_name" env:"MS_DCR_STREAM" valid:"minstringlength(3)"`
		} `yaml:"dcr"`

		ResourceGroup string `yaml:"resource_group" env:"MS_RSG_ID" valid:"minstringlength(3)"`
		WorkspaceName string `yaml:"workspace_name" env:"MS_WS_NAME" valid:"minstringlength(3)"`
	} `yaml:"microsoft"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if c.Notion.Lookback == "" {
		c.Notion.Lookback = defaultLookback
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("could not load environment: %v", err)
	}

	return nil
}
