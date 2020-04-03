package context

import (
	//	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Appname string
	Version string

	Db *Mongodbconfig
}

type Mongodbconfig struct {
	ConnectionString string
	User             string
	Password         string
	Dbname           string
	CollectionName   string
}

func LoadConfigMgo(path string) (Config, error) {
	configs := viper.New()
	configs.SetConfigName("Config")
	configs.AddConfigPath(path)

	err := configs.ReadInConfig()
	if err != nil {
		log.Println("error", err)
		return Config{}, err
	}
	//	log.Println(configs.Get("User"))
	//log.Println(configs)

	var config Config
	err = configs.Unmarshal(&config)
	log.Println("configll", config.Db.User)

	return config, err
}
