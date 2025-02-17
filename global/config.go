package global

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func SetConfigName(configname string) {
	Config.configName = configname
}

func SetConfigType(typename string) {
	Config.configType = typename
}

func AddConfigPath(path string) {
	Config.configPath = append(Config.configPath, path)
}

// Config Key
const (

	//
	GATAWAY_NAME     = "gateway.name"
	GATAWAY_SN       = "gateway.sn"
	SYS_ENV          = "system.environment"
	SYS_WATCH_CONFIG = "system.watchConfig"
	LOG_LEVEL        = "logConf.level"
	LOG_PATH         = "logConf.path"
	LOG_PARTTEN      = "logConf.partten"
	LOG_MAXAGE       = "logConf.maxAge"
	LOG_ROTATION     = "logConf.rotationTime"
	LOG_USE_COMPRESS = "logConf.compress"
	//
	FANUC_LIB = "fanuc.lib"
	FANUC_LOG = "fanuc.log"
	//
	CORE_SN         = "core.sn"
	CORE_WORKPATH   = "core.workpath"
	CORE_PLUGINPATH = "core.pluginpath"
	//comm
	//mqtt
	MQTT_BROKER = "comm.mqtt.broker"
	MQTT_USER   = "comm.mqtt.user"
	MQTT_PASS   = "comm.mqtt.pass"
	MQTT_TOPIC  = "comm.mqtt.topic"
)

//setting defaults
//reading from environment variables

//Viper uses the following precedence order. Each item takes precedence over the item below it:

// explicit call to Set
// flag
// env
// config
// key/value store
// default
func init() {
	Config = &config{
		configName: "config",
		configType: "yaml",
		configPath: []string{"."},
	}
}

var Config *config

type config struct {
	configName string
	configType string
	configPath []string
	GateWay    gateway `yaml:"gateway"`
	System     system  `yaml:"system"`
	LogConf    logConf `yaml:"logConf"`
}

type gateway struct {
	Name string `yaml:"name"`
	SN   string `yaml:"sn"`
}
type system struct {
	Environment string `yaml:"environment"`
}

type logConf struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	Partten  string `yaml:"partten"`
	MaxAge   string `yaml:"maxAge"`
	Rotation string `yaml:"rotationTime"`
	Compress bool   `yaml:"compress"`
}

func LoadConfig(watch bool) {
	loadConfig()
	if watch {
		//Watching and re-reading config files
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e)
			fmt.Println("Config file changed:", e.Name)
			loadConfig()
		})
		viper.WatchConfig()
	}
}

func loadConfig() {
	viper.SetConfigName(Config.configName)
	viper.SetConfigType(Config.configType)
	for _, path := range Config.configPath {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(fmt.Errorf("fatal error config file: %w", err))
		} else {
			log.Printf("config file was found but another error was produced: %v\n", err)
		}
	}
	// Writing Config Files
	// viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	// viper.SafeWriteConfig()
	// viper.WriteConfigAs("/path/to/my/.config")
	// viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
	// viper.SafeWriteConfigAs("/path/to/my/.other_config")
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
	PreLog("配置文件", zap.Any("Config", viper.AllSettings()))
}
