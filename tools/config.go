package tools

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/go-yaml/yaml"
)

var config *Config

func loadConfigYaml(p string) ([]byte, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Config struct {
	SlackAPI SlackApiConfig `yaml:"slackApi"`
}
type SlackApiConfig struct {
	Messaging MessagingConfig `yaml:"messaging"`
}
type MessagingConfig struct {
	WebhookURL string `yaml:"webhookUrl"`
}

func GetConfig() (*Config, error) {
	wd, _ := os.Getwd()
	log.Printf("wd=%v", wd)

	if config != nil {
		return config, nil
	}

	config := &Config{}
	configPath := path.Join(wd, "/config.yaml")
	bytes, err := loadConfigYaml(configPath)
	if err != nil {
		configPath = path.Join(wd, "../config.yaml")
		bytes, err = loadConfigYaml(configPath)
		if err != nil {
			return nil, err
		}
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
