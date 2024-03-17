package toes

import (
	"github.com/starjun/gotools/v2"
	"log"
	"os"
	"path"
	"strings"
)

var global_conf_Tpl = `
package global

import "time"

type Config struct {
	Log         Log         {{at}}mapstructure:"log" json:"log" yaml:"log"{{at}}
	Seckey      Seckey      {{at}}mapstructure:"seckey" json:"seckey" yaml:"seckey"{{at}}
	CheckHeader CheckHeader {{at}}mapstructure:"checkHeader" json:"checkHeader" yaml:"checkHeader"{{at}}
	Server      Server      {{at}}mapstructure:"server" json:"server" yaml:"server"{{at}}
	Tls         Tls         {{at}}mapstructure:"tls" json:"tls" yaml:"tls"{{at}}
	Mysql       Mysql       {{at}}mapstructure:"mysql" json:"mysql" yaml:"mysql"{{at}}
	Redis       Redis       {{at}}mapstructure:"redis" json:"redis" yaml:"redis"{{at}}
	Header      Header      {{at}}mapstructure:"header" json:"header" yaml:"header"{{at}}
}

type Tls struct {
	Addr string {{at}}mapstructure:"addr" json:"addr" yaml:"addr"{{at}}
	Cert string {{at}}mapstructure:"cert" json:"cert" yaml:"cert"{{at}}
	Key  string {{at}}mapstructure:"key" json:"key" yaml:"key"{{at}}
}

type Log struct {
	Format  string {{at}}mapstructure:"format" json:"format" yaml:"format"{{at}}
	Console bool   {{at}}mapstructure:"console" json:"console" yaml:"console"{{at}}
	Path    string {{at}}mapstructure:"path" json:"path" yaml:"path"{{at}}
	Level   string {{at}}mapstructure:"level" json:"level" yaml:"level"{{at}}
	Days    int    {{at}}mapstructure:"days" json:"days" yaml:"days"{{at}}
}

type Seckey struct {
	Basekey    string {{at}}mapstructure:"basekey" json:"basekey" yaml:"basekey"{{at}}
	Jwtttl     int    {{at}}mapstructure:"jwtttl" json:"jwtttl" yaml:"jwtttl"{{at}}
	Pproftoken string {{at}}mapstructure:"pproftoken" json:"pproftoken" yaml:"pproftoken"{{at}}
}

type CheckHeader struct {
	Nonce             bool    {{at}}mapstructure:"nonce" json:"nonce" yaml:"nonce"{{at}}
	NonceCacheSeconds int     {{at}}mapstructure:"nonceCacheSeconds" json:"nonceCacheSeconds" yaml:"nonceCacheSeconds"{{at}}
	Time              bool    {{at}}mapstructure:"time" json:"time" yaml:"time"{{at}}
	Seconds           float64 {{at}}mapstructure:"seconds" json:"seconds" yaml:"seconds"{{at}}
	Sign              bool    {{at}}mapstructure:"sign" json:"sign" yaml:"sign"{{at}}
	All               bool    {{at}}mapstructure:"all" json:"all" yaml:"all"{{at}}
}

type Server struct {
	Mode string {{at}}mapstructure:"mode" json:"mode" yaml:"mode"{{at}}
	Addr string {{at}}mapstructure:"addr" json:"addr" yaml:"addr"{{at}}
}

type Mysql struct {
	Host                  string        {{at}}mapstructure:"host" json:"host" yaml:"host"{{at}}
	Username              string        {{at}}mapstructure:"username" json:"username" yaml:"username"{{at}}
	Password              string        {{at}}mapstructure:"password" json:"password" yaml:"password"{{at}}
	MaxOpenConnections    int           {{at}}mapstructure:"maxOpenConnections" json:"maxOpenConnections" yaml:"maxOpenConnections"{{at}}
	MaxConnectionLifeTime time.Duration {{at}}mapstructure:"maxConnectionLifeTime" json:"maxConnectionLifeTime" yaml:"maxConnectionLifeTime"{{at}}
	LogLevel              int           {{at}}mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"{{at}}
	PasswordMode          string        {{at}}mapstructure:"passwordMode" json:"passwordMode" yaml:"passwordMode"{{at}}
	Database              string        {{at}}mapstructure:"database" json:"database" yaml:"database"{{at}}
	MaxIdleConnections    int           {{at}}mapstructure:"maxIdleConnections" json:"maxIdleConnections" yaml:"maxIdleConnections"{{at}}
}

type Redis struct {
	Password string {{at}}mapstructure:"password" json:"password" yaml:"password"{{at}}
	Host     string {{at}}mapstructure:"host" json:"host" yaml:"host"{{at}}
	Username string {{at}}mapstructure:"username" json:"username" yaml:"username"{{at}}
}

type Header struct {
	Realip    string {{at}}mapstructure:"realip" json:"realip" yaml:"realip"{{at}}
	Requestid string {{at}}mapstructure:"requestid" json:"requestid" yaml:"requestid"{{at}}
}

type EnvCfg struct {
	MyName string
	MyId   string
}

`

func Global_ConfFileGen(_Path string) {
	filepath := path.Join(_Path, "config.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return
	}
	if re {
		log.Println(filepath, "已存在")
		return
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
	}
	global_conf_Tpl = strings.Replace(global_conf_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_conf_Tpl)
	f.Close()
}
