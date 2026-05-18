package config

//引入viper: go get -u github.com/spf13/viper
import (
	"log"

	"github.com/spf13/viper"
)

// 参考config.yaml文件的结构定义Config结构体
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
}

//引入gin: go get  -u github.com/gin-gonic/gin

// 全局变量AppConfig用于存储加载的配置
var AppConfig *Config

func InitConfig() {
	//创建viper实例
	viper.SetConfigName("config")   //配置文件名,不需要后缀
	viper.SetConfigType("yaml")     //配置文件类型,yml或json等
	viper.AddConfigPath("./config") //配置文件路径
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	InitDB()
	InitRedis()
}
