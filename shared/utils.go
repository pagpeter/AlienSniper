package shared

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetIP() string {
	url := "https://api.ipify.org"

	resp, err := http.Get(url)
	if err != nil {
		// log.Fatalln(err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// log.Fatalln(err)
		return ""
	}
	return string(body)
}

func Input(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func WriteFile(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0644)
}

func IsInMap(key string, m map[string]int) bool {
	_, ok := m[key]
	return ok
}
