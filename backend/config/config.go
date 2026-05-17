package config

//引入viper: go get -u github.com/spf13/viper
import "github.com/spf13/viper"

//参考config.yaml文件的结构定义Config结构体
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

//引入gin: go get  -u github.com/gin-gonic/gin

//全局变量AppConfig用于存储加载的配置
var AppConfig *Config

func InitConfig() {
	//创建viper实例
	viper.SetConfigName("config") //配置文件名,不需要后缀

}
