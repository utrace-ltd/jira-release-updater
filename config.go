package main

import "github.com/spf13/viper"

type Config struct {
	Jira struct{
		URL string
		User string
		Password string
		Project struct{
			Id string
		}
		Issue struct{
			Pattern string
		}
	}
}

func NewConfig() (*Config, error)  {
	config := &Config{}
	viper.SetDefault("jira.url", "https://utrace.atlassian.net")
	viper.SetDefault("jira.user", "user@utrace.ru")
	viper.SetDefault("jira.password", "secret")
	viper.SetDefault("jira.project.id", "10001")
	viper.SetDefault("jira.issue.pattern", "(ISSUE-+.)")

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/jira-release-updater/")
	viper.AddConfigPath("$HOME/.jira-release-updater/")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
