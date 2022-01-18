package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var clear map[string]func() //create a map for storing clear funcs

// https://stackoverflow.com/a/22896706
func init() {
	clear = make(map[string]func()) //Initialize it
	unix := func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = unix
	clear["darwin"] = unix

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! Supported: Darwin (MacOS), Linux, Windows")
	}
}

func main() {
	ClearScreen()
	PrintLogo()
	PrintStats()
	fmt.Println("\n")
	time.Sleep(time.Second * 100)
}

func PrintLogo() {
	text := "               _,--=--._\n"
	text += "             ,'    _    `.\n"
	text += "            -    _(_)_o   - \n"
	text += "       ____'    /_  _/]    `____\n"
	text += "-=====::(+):::::::::::::::::(+)::=====-\n"
	text += `         (+).""""""""""""",(+)` + "\n"
	text += "             .           ,\n"
	text += "               `  -=-  '\n"
	text += "\n\n"
	text += "Alien: Faster than every cowboy\n\n"
	fmt.Println(text)
}

func PrintStats() {
	fmt.Println("Running as host\n")
	fmt.Println("Servers connected: " + "7")        //strconv.Itoa(len(servers)))
	fmt.Println("Names attempted to snipe: " + "5") //strconv.Itoa(len(names)))
	fmt.Println("Names successfully sniped: " + "3")
}
