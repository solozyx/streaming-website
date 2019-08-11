package conf

import (
	"encoding/json"
	"io/ioutil"
)

var (
	G_config *Config
)

type Config struct {
	ApiPort         int    `json:"apiPort"`
	ApiReadTimeout  int    `json:"apiReadTimeout"`
	ApiWriteTimeout int    `json:"apiWriteTimeout"`

	ApiUserRegister string `json:"userRegister"`
	ApiUserLogin    string `json:"userLogin"`
}

func InitConfig(confFilePath string) (err error) {
	var (
		content []byte
		config  Config
	)
	if content, err = ioutil.ReadFile(confFilePath); err != nil {
		return
	}
	if err = json.Unmarshal(content, &config); err != nil {
		return
	}
	G_config = &config
	return
}
