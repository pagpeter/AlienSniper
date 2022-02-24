package types

import (
	utils "Alien/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

type Requests struct {
	Giftcard  int `json:"giftcard"`
	Microsoft int `json:"microsoft"`
}

type TLS struct {
	Active bool   `json:"active"`
	Key    string `json:"key"`
	Cert   string `json:"cert"`
}

type Config struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Requests Requests `json:"requests"`
	Token    string   `json:"token"`
	TLS      TLS      `json:"tls"`
}

func (c *Config) LoadFromFile() {
	// Load a config file

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalln("Failed to open config file: ", err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func (c *Config) SaveToFile() {
	// Save the config to a file

	json, _ := json.MarshalIndent(c, "", "  ")
	ioutil.WriteFile("config.json", json, 0644)
}

type IP struct {
	Query string
}

func Configure() *Config {
	// Start the configuring process, used to generate a config file

	c := Config{}

	if runtime.GOOS == "windows" {
		c.Host = "127.0.0.1"
	} else {
		c.Host = "0.0.0.0"
	}

	fmt.Println("Welcome, we are now configuring Alien for you. ")

	c.Token = utils.GenerateToken(20)
	fmt.Println("Generated token: ", c.Token)
	fmt.Println("Using host at: " + c.Host)

	fmt.Print("Port: ")
	fmt.Scan(&c.Port)

	fmt.Printf("Using port: %v\n", c.Port)

	c.Requests.Giftcard = 2
	fmt.Println("Using 2 giftcard requests")

	c.Requests.Microsoft = 2
	fmt.Println("Using 2 microsoft requests")

	return &c
}
