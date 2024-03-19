package toes

import (
	"errors"
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
    Jwtkey     int    {{at}}mapstructure:"jwtkey" json:"jwtkey" yaml:"jwtkey"{{at}}
	Jwtttl     int    {{at}}mapstructure:"jwtttl" json:"jwtttl" yaml:"jwtttl"{{at}}
	Pproftoken string {{at}}mapstructure:"pproftoken" json:"pproftoken" yaml:"pproftoken"{{at}}
}

type CheckHeader struct {
	Nonce             bool    {{at}}mapstructure:"nonce" json:"nonce" yaml:"nonce"{{at}}
	NonceCacheSeconds int     {{at}}mapstructure:"nonceCacheSeconds" json:"nonceCacheSeconds" yaml:"nonceCacheSeconds"{{at}}
	Time              bool    {{at}}mapstructure:"time" json:"time" yaml:"time"{{at}}
	Seconds           float64 {{at}}mapstructure:"seconds" json:"seconds" yaml:"seconds"{{at}}
	Sign              bool    {{at}}mapstructure:"sign" json:"sign" yaml:"sign"{{at}}
    Key               bool    {{at}}mapstructure:"key" json:"key" yaml:"key"{{at}}
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

var global_db_Tpl = `
package global

import (
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysqllogger "gorm.io/gorm/logger"
	"strings"
)

var (
	DB *gorm.DB
)

// InitStore 读取 db 配置，创建 gorm.DB 实例，并初始化 store 层.
func InitStore() error {
	var err error
	DB, err = newMySQL(&Cfg.Mysql)
	if err != nil {
		LogErrorw("mysql连接失败", "Subject", "mysql", "Result", err)
		cobra.CheckErr(err)
	}

	LogDebugw("init db success")
	return nil
}

// DSN returns mysql dsn.
func (o *Mysql) DSN() string {
	return fmt.Sprintf({{at}}%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s{{at}},
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local",
	)
}

// newMySQL create a new gorm db instance with the given options.
func newMySQL(opts *Mysql) (*gorm.DB, error) {
	logLevel := mysqllogger.Silent
	if opts.LogLevel != 0 {
		logLevel = mysqllogger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: mysqllogger.Default.LogMode(logLevel),
		//SkipDefaultTransaction: true,
		//PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

`

var global_gl_Tpl = `
package global

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// RecommendedName defines the default project name.
	RecommendedName = "toes"
	LogTmFmt        = "2006-01-02 15:04:05"
)

var (
	CfgFile     string
	Cache       *cache.Cache
	RedisClient *redis.Client
	Ctx         = context.Background() // redis 使用的
	Cfg         *Config
)

var (
	mu          sync.Mutex
	defaultName = RecommendedName + ".config.yaml"

	// RecommendedEnvPrefix defines the ENV prefix used by all service.
	RecommendedEnvPrefix = strings.ToUpper(RecommendedName)
)

// 最先进行初始化的
func InitConfig() {
	initConfig(CfgFile)
}

func initConfig(cfgpath string) {
	if cfgpath != "" {
		viper.SetConfigFile(cfgpath)
	} else {
		// 获取目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Add {{at}}$HOME/<RecommendedHomeDir>{{at}} & {{at}}.{{at}}
		viper.AddConfigPath(filepath.Join(home, "conf"))
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(".", "conf"))

		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read viper configuration file", "err", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Println("config unmarshal err", "err", err)
	}

	// Print using config file.
	log.Println("Using config file", "file", viper.ConfigFileUsed())

	// Watch config file
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file updated")
		if err := viper.Unmarshal(&Cfg); err != nil {
			log.Println("config unmarshal err", "err", err)

			return
		}

		// 暂时日志不更新
	})
}

func InitLocalCache() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
	//Cache.Set("foo", "bar", cache.DefaultExpiration)
}

func InitRedis() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Host,
		Password: Cfg.Redis.Password,
		Username: Cfg.Redis.Username,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("RedisClient.Ping error")
		cobra.CheckErr(err)
	}
}

`

var global_l2cache_Tpl = `
package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	L2cache l2cache
)

type l2cache struct {
	Lcache     *cache.Cache
	Ltime      int
	RedisCache *redis.Client
	Rtime      int
	Pre        string
}

func InitL2cache(pre string, ltime, rtime int) {
	if ltime == 0 {
		ltime = 100 // 本地 cache 默认100s
	}
	if rtime == 0 {
		rtime = 200 // redis 缓存 200s
	}
	L2cache.RedisCache = RedisClient
	L2cache.Lcache = Cache
}

func (l *l2cache) Set(key, value string) {
	l.Lcache.Set(l.Pre+key, value, time.Duration(l.Ltime)*time.Second)
	l.RedisCache.Set(Ctx, l.Pre+key, value, time.Duration(l.Rtime)*time.Second)
}

func (l *l2cache) Get(key string) (value string, err error) {
	if x, found := l.Lcache.Get(l.Pre + "foo"); found {
		value = x.(string)
		return
	}
	value, err = l.RedisCache.Get(Ctx, l.Pre+key).Result()
	if err != nil {
		l.Lcache.Set(l.Pre+key, value, time.Duration(l.Ltime)*time.Second)
	}
	return
}

`

var global_log_Tpl = `
package global

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	logger *zap.Logger
)

func InitLog(_log *Log) {
	mu.Lock()
	defer mu.Unlock()
	getLogger(_log)
}

func getLogger(_log *Log) {
	var level zapcore.Level
	switch _log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller_line"
	encoderConfig.FunctionKey = zapcore.OmitKey
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	//encoderConfig.LineEnding = "\r"
	encoderConfig.EncodeLevel = cEncodeLevel
	//encoderConfig.EncodeTime = cEncodeTime
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = cEncodeCaller

	/*
		LineEnding:     zapcore.DefaultLineEnding,     //输出的分割符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //序列化字符串的大小写
		//EncodeTime:          zapcore.ISO8601TimeEncoder,     //时间的编码格式
		EncodeTime:          EncodeTime,                     //时间自定义的
		EncodeDuration:      zapcore.SecondsDurationEncoder, //时间显示的位数
		EncodeCaller:        zapcore.ShortCallerEncoder,     //输出的运行文件路径长度
		EncodeName:          zapcore.FullNameEncoder,        //可选的
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "", //控制台格式时，每个字段间的分割符,不配置默认即可
	*/
	var Encoder zapcore.Encoder
	if _log.Format == "console" {
		Encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		Encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	WriteSyncer := &lumberjack.Logger{
		Filename:   _log.Path,
		MaxSize:    300,
		MaxBackups: 3,
		MaxAge:     _log.Days,
	}

	writes := []zapcore.WriteSyncer{zapcore.AddSync(WriteSyncer)}
	if _log.Console {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(Encoder,
		zapcore.NewMultiWriteSyncer(writes...),
		level)

	//6.构造日志
	//设置为开发模式会记录panic
	development := zap.Development()
	//caller := zap.WithCaller(true)
	//构造一个字段
	//zap.Fields(zap.String("appName", "demozap"))
	//通过传入的配置实例化一个日志
	logger = zap.New(core, development, zap.AddCaller())

	// 替换全局 zap log
	zap.ReplaceGlobals(logger)
	// 全局使用 eg
	// zap.S().Info("hello")

	// 把标准库的 log.Logger 的 info 级别的输出重定向到 zap.Logger
	zap.RedirectStdLog(logger)
}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(LogTmFmt) + "]")
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

func LogSync() {
	logger.Sync()
}

// Debugw 输出 debug 级别的日志.
func LogDebugw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Debugw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Infow 输出 info 级别的日志.
func LogInfow(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Infow(msg, keysAndValues...)
	//defer logger.Sync()
}

// Warnw 输出 warning 级别的日志.
func LogWarnw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Warnw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Errorw 输出 error 级别的日志.
func LogErrorw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Errorw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Panicw 输出 panic 级别的日志.
func LogPanicw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Panicw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Fatalw 输出 fatal 级别的日志.
func LogFatalw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Fatalw(msg, keysAndValues...)
	//defer logger.Sync()
}

// -- web中间件记录日志
func LogGin(ctx context.Context) *zap.Logger {

	_logger := logger.With(zap.Any("WEB", "GIN"))
	if requestID := ctx.Value(Cfg.Header.Requestid); requestID != nil {
		_logger = _logger.With(zap.Any("Traceid", requestID))
	}
	return _logger
}

`

func (b *ToesGen) Global_ConfFileGen() error {
	filepath := path.Join(b.globalPath, "config.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
		return e
	}
	global_conf_Tpl = strings.Replace(global_conf_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_conf_Tpl)
	f.Close()
	return nil
}

func (b *ToesGen) Global_DbFileGen() error {
	filepath := path.Join(b.globalPath, "db.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
	}
	global_db_Tpl = strings.Replace(global_db_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_db_Tpl)
	f.Close()
	return nil
}

func (b *ToesGen) Global_GlFileGen() error {
	filepath := path.Join(b.globalPath, "gl.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
	}
	global_gl_Tpl = strings.Replace(global_gl_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_gl_Tpl)
	f.Close()
	return nil
}

func (b *ToesGen) Global_L2cacheFileGen() error {
	filepath := path.Join(b.globalPath, "l2cache.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
	}
	global_l2cache_Tpl = strings.Replace(global_l2cache_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_l2cache_Tpl)
	f.Close()
	return nil
}

func (b *ToesGen) Global_LogFileGen() error {
	filepath := path.Join(b.globalPath, "log.go")
	re, err := gotools.PathExists(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	if re {
		log.Println(filepath, "已存在")
		return errors.New(filepath + " 已存在")
	}

	f, e := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)
	if e != nil {
		log.Println(e)
	}
	global_log_Tpl = strings.Replace(global_log_Tpl, "{{at}}", "`", -1)
	f.WriteString(global_log_Tpl)
	f.Close()
	return nil
}
