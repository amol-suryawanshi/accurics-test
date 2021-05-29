package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	//FullAccessScope represents full access for a repo
	FullAccessScope = "fullAccess"
)

//ServerConfig is in-memory configuration
var ServerConfig *configuration

type configuration struct {
	OAuthRedirectURL string `json:"redirectURL"`
	Scope            string `json:"scope"`
	GithubOAuthURL   string `json:"gitOAuthURL"`
	AllowSignUp      string `json:"allowSignup"`
	GithubAPIURL     string `json:"gitAPIhost"`
	ClientID         string `json:"clientID"`
	ClientSecret     string `json:"clientSecret"`
}

//ReadConfig reads configuration from given path
func ReadConfig(configFilePath string) error {
	ServerConfig = &configuration{}
	byteArr, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteArr, ServerConfig)
	if err != nil {
		return err
	}
	return nil
}

//LoadApplication loads config and logger
func LoadApplication(configFilePath string) error {
	err := ReadConfig(configFilePath)
	if err != nil {
		return err
	}
	return nil
}
