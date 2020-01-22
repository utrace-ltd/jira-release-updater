package main

import "github.com/spf13/viper"

type Config struct {
	Jira struct {
		URL      string
		User     string
		Password string
		Project  struct {
			Id string
		}
		Issue struct {
			Pattern string
		}
	}
}

func NewConfig() (*Config, error) {
	config := &Config{}

	v := viper.NewWithOptions(viper.KeyDelimiter("_"))

	v.SetDefault("jira_url", "https://utrace.atlassian.net")
	v.SetDefault("jira_user", "user@utrace.ru")
	v.SetDefault("jira_password", "secret")
	v.SetDefault("jira_project_id", "10001")
	v.SetDefault("jira_issue_pattern", "(ISSUE-+.)")

	v.AutomaticEnv()

	v.SetConfigName("config")
	v.AddConfigPath("/etc/jira-release-updater/")
	v.AddConfigPath("$HOME/.jira-release-updater/")
	v.AddConfigPath(".")

	_ = v.ReadInConfig()

	err := v.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
