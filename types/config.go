package types

import (
	utils "Alien/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Requests struct {
	Giftcard  int `json:"giftcard"`
	Mojang    int `json:"mojang"`
	Microsoft int `json:"microsoft"`
}

type Config struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Requests Requests `json:"requests"`
	Token    string   `json:"token"`
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

func Configure() Config {
	// Start the configuring process, used to generate a config file

	IP := utils.GetIP()
	c := Config{}

	fmt.Println("Welcome, we are now configuring Alien for you. ")
	fmt.Println("Press enter to use the default values.\n")

	c.Token = utils.GenerateToken(20)
	fmt.Println("Generated token: ", c.Token)

	c.Host = utils.Input(fmt.Sprintf("IP address of host API (%s):\n> ", IP))
	if c.Host == "" {
		c.Host = IP
		fmt.Println("\nUsing host at: " + c.Host)
	}

	c.Port = utils.ToInt(utils.Input(fmt.Sprintf("Port of host API (20514):\n> ")))
	if c.Port == 0 {
		c.Port = 20514
		fmt.Println("\nUsing port 20514")
	}

	c.Requests.Giftcard = utils.ToInt(utils.Input(fmt.Sprintf("Giftcard requests (2):\n> ")))
	if c.Requests.Giftcard == 0 {
		c.Requests.Giftcard = 2
		fmt.Println("\nUsing 2 giftcard requests")
	}

	c.Requests.Mojang = utils.ToInt(utils.Input(fmt.Sprintf("Mojang requests (2):\n> ")))
	if c.Requests.Mojang == 0 {
		c.Requests.Mojang = 2
		fmt.Println("\nUsing 2 mojang requests")
	}

	c.Requests.Microsoft = utils.ToInt(utils.Input(fmt.Sprintf("Microsoft requests (2):\n> ")))
	if c.Requests.Microsoft == 0 {
		c.Requests.Microsoft = 2
		fmt.Println("\nUsing 2 microsoft requests")
	}

	return c
}
