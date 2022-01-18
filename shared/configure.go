package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c *Config) LoadFromFile() {
	// Load a config file

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func (c *Config) SaveToFile() {
	// Save the config to a file

	json, _ := json.MarshalIndent(c, "", "  ")
	err := ioutil.WriteFile("config.json", json, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func Configure() Config {
	// Start the configuring process, used to generate a config file

	IP := GetIP()
	c := Config{}

	fmt.Println("Welcome, we are now configuring Alien for you. ")
	fmt.Println("Press enter to use the default values.\n")

	c.Host = Input(fmt.Sprintf("IP address of host API (%s):\n> ", IP))
	if c.Host == "" {
		c.Host = IP
		fmt.Println("Using host at: " + c.Host)
	}

	fmt.Println("\n")

	c.Port = ToInt(Input(fmt.Sprintf("Port of host API (20514):\n> ")))
	if c.Port == 0 {
		c.Port = 20514
		fmt.Println("Using port 20514")
	}

	return c
}
