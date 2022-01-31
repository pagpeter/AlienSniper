package types

import (
	utils "Alien/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Requests struct {
	Giftcard  int `json:"giftcard"`
	Mojang    int `json:"mojang"`
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

func getip2() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func Configure() *Config {
	// Start the configuring process, used to generate a config file

	// IP := utils.GetIP()
	c := Config{}

	fmt.Println("Welcome, we are now configuring Alien for you. ")

	c.Token = utils.GenerateToken(20)
	fmt.Println("Generated token: ", c.Token)

	c.Host = getip2()
	fmt.Println("\nUsing host at: " + c.Host)

	c.Port = 21615
	fmt.Println("\nUsing port 21615")

	c.Requests.Giftcard = 2
	fmt.Println("\nUsing 2 giftcard requests")

	c.Requests.Mojang = 2
	fmt.Println("\nUsing 2 mojang requests")

	c.Requests.Microsoft = 2
	fmt.Println("\nUsing 2 microsoft requests")

	return &c
}
