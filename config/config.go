package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/starjun/gotools/v2"
	"github.com/starjun/toes-cli/util"
	"gorm.io/gorm"
)

var (
	Cfg               *Config
	DB                *gorm.DB
	Dsn               string
	CfgPath           string
	Style             string
	Table             string
	PackageName       string
	ToesRoot          string
	CreateQueryConfig bool
	Version           string = "v1"
	Path              string
	RootPackage       string
)

type Config struct {
	Version           string       `mapstructure:"version" json:"version" yaml:"version"`
	Style             string       `mapstructure:"style" json:"style" yaml:"style"`
	PackageName       string       `mapstructure:"packageName" json:"packageName" yaml:"packageName"`
	ToesRoot          string       `mapstructure:"toesRoot" json:"toesRoot" yaml:"toesRoot"`
	CreateQueryConfig bool         `mapstructure:"createQueryConfig" json:"createQueryConfig" yaml:"createQueryConfig"`
	MysqlOptions      MySQLOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

type MySQLOptions struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
}

func (o *MySQLOptions) DSN() string {
	// root:mima@tcp(127.0.0.1)/mydb?charset=utf8&parseTime=true&loc=Local
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local",
	)
}

func GetPath(tableName string) string {
	if Style != "toes" {
		return tableName
	}
	if strings.TrimSpace(Table) == "" || strings.TrimSpace(Table) == "*" {
		return tableName
	}
	if strings.TrimSpace(Path) != "" {
		return Path
	}

	return tableName
}

func GetDSN() string {
	if Cfg != nil {
		return strings.TrimSpace(Cfg.MysqlOptions.DSN())
	}

	return strings.TrimSpace(Dsn)
}

func GetRootPackage() string {
	if Cfg != nil {
		return Cfg.PackageName
	}

	return PackageName
}

func GetCreateQueryConfig() bool {
	if Cfg != nil {
		return Cfg.CreateQueryConfig
	}

	return CreateQueryConfig
}
func GetVersion() string {
	if Cfg != nil {
		return Cfg.Version
	}

	return Version
}

func NewDefaultConfig() *Config {
	return &Config{
		Version:           "v1",
		PackageName:       "toesdemo",
		ToesRoot:          "./",
		CreateQueryConfig: false,
		MysqlOptions: MySQLOptions{
			Host:     "127.0.0.1",
			Username: "root",
			Password: "",
			Database: "test",
		},
	}
}

func LoadConfig(cfg string) {
	re, _ := gotools.PathExists(cfg)
	if !re {
		createConfig()
	}
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.toes")
	}

	if err := viper.ReadInConfig(); err != nil {

		return
	}

	Cfg = NewDefaultConfig()
	if err := viper.Unmarshal(Cfg); err != nil {
		fmt.Println(err)
	}
	//Cfg.MysqlOptions.Password = util.DecryptString(Cfg.MysqlOptions.Password)
}

const (
	tmplConfig = `
mysql:
  host: 127.0.0.1 # MySQL 机器 ip 和端口，默认 127.0.0.1:3306
  username: root # MySQL 用户名(建议授权最小权限集)
  passWord: 123 # MySQL 用户密码
  database: test # 系统所用的数据库名

version:     "v1.0.0"
PackageName: "library"
toesRoot: "./"
createQueryConfig: false
`
)

func createConfig() {
	fileName := "./config.toes.yaml"
	re, _ := gotools.PathExists(fileName)
	if !re {
		log.Printf("'%s' not exists. Do you want to create it? [Yes|No]\n", fileName)
		if util.AskForConfirmation() {
			ioutil.WriteFile(fileName, []byte(tmplConfig), 0777)
		}
	}
}
