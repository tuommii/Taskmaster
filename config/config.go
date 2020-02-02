package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config ...
type Config struct {
	Entries map[string]*Entry
}

// Entry represents one key
type Entry struct {
	Autorestart bool   `json:"autorestart"`
	Autostart   bool   `json:"autostart"`
	Command     string `json:"command"`
	StdErrLog   string `json:"stderr_log"`
	StdOutLog   string `json:"stdout_log"`
}

// LoadConfig ...
func LoadConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error while opening config file: ", err)
		panic(err)
	}
	config := &Config{}
	err = json.Unmarshal([]byte(file), &config.Entries)
	return config
}

// Get ...
func (c Config) Get(name string) *Entry {
	for k := range c.Entries {
		if k == name {
			return c.Entries[name]
		}
	}
	return nil
}

// func (c ConfigFile) Print() {
// 	for k, v := range c {
// 		fmt.Printf("key[%s] value[%v]", k, v)
// 	}
// }
