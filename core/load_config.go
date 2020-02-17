package core

import (
	"github.com/TRON-US/btfs-migration-toolkit/conf"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf *conf.Config

//init config file to viper
func InitConfig(cfg string) error {
	if err := initViper(cfg); err != nil {
		return err
	}
	return nil
}

//init config file to viper
func initViper(cfg string) error {
	if cfg != "" {
		viper.SetConfigFile(cfg) //  set config file
	} else {
		viper.AddConfigPath("conf") // set default file
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml") // set config file to yaml file
	viper.AutomaticEnv()        // read env variable

	if err := getNewConfig(); err != nil {
		return err
	}

	watchConfig()
	return nil
}

func getNewConfig() error {
	var err error
	if err = viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}
	Conf = &conf.Config{}

	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}

	return nil
}

// watch file changing
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		_ = getNewConfig()
	})
}
