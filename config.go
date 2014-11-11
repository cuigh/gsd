package gsd

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

var (
	ConfigPath string
	_Configs   map[string]*Config
	_Locker    sync.Mutex
)

type SettingMap map[string]string

func (this SettingMap) String(key string, defaultValue string) string {
	v, ok := this[key]
	if ok {
		return v
	}

	return defaultValue
}

func (this SettingMap) Int(key string, defaultValue int) int {
	v, ok := this[key]
	if ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}

	return defaultValue
}

// database settings
type Config struct {
	Name     string
	Provider string
	Driver   string
	Settings SettingMap
}

// GetConfig return configuration of specific database in database.sql.conf file
func GetConfig(name string) (cfg *Config, err error) {
	if _Configs == nil {
		_Locker.Lock()
		defer _Locker.Unlock()

		if _Configs == nil {
			if ConfigPath == "" {
				err = fmt.Errorf("You must set [ConfigPath] first for locating databases.")
			}

			var file *os.File
			file, err = os.Open(ConfigPath)
			if err != nil {
				return
			}
			defer file.Close()

			decoder := xml.NewDecoder(file)
			_Configs, err = loadConfig(decoder)
			if err != nil {
				return
			}
		}
	}

	var ok bool
	cfg, ok = _Configs[name]
	if !ok {
		err = fmt.Errorf("cannot find the configuration of database [%s]", name)
	}
	return
}

func loadConfig(decoder *xml.Decoder) (configs map[string]*Config, err error) {
	configs = make(map[string]*Config)

	var (
		t     xml.Token
		cfg   *Config
		name  string
		value string
	)

	for {
		t, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		switch token := t.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "database":
				cfg = &Config{Settings: SettingMap{}}
				for _, attr := range token.Attr {
					switch attr.Name.Local {
					case "name":
						cfg.Name = attr.Value
					case "provider":
						cfg.Provider = attr.Value
					case "driver":
						cfg.Driver = attr.Value
					}
				}
			case "setting":
				for _, attr := range token.Attr {
					switch attr.Name.Local {
					case "name":
						name = attr.Value
					case "value":
						value = attr.Value
					}
				}
				cfg.Settings[name] = value
			}
		case xml.EndElement:
			if token.Name.Local == "database" {
				if cfg.Driver == "" {
					cfg.Driver = cfg.Provider
				}
				configs[cfg.Name] = cfg
			}
		}
	}
	return
}
