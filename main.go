package main

import (
	"Alien/host"
	shared "Alien/shared"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const version = "0.0.1"

var clear map[string]func() //create a map for storing clear funcs
var prefix = map[string]string{
	"windows": "",
	"darwin":  "./",
	"linux":   "./",
}

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
	// Clear screen on supported platforms

	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! Supported: Darwin (MacOS), Linux, Windows")
	}
}

func main() {
	HandleArgs()
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

func HandleArgs() {
	// Handle command line arguments

	if len(os.Args) > 1 {
		prefix := prefix[runtime.GOOS]
		usage := fmt.Sprintf("Usage: %salien [options]\n\n", prefix)
		usage += "Options:\n"
		usage += "    help: Print this help message\n"
		usage += "    version: Print the version number\n"
		usage += "    configure: Configure the application\n"
		usage += "    start: Start the CLI\n"
		usage += "    host: Start as the host in the background\n"
		usage += "    node: Start as a node in the background\n"

		if len(os.Args) == 1 {
			fmt.Println(usage)
			os.Exit(0)
		}

		arg := os.Args[1]
		switch arg {
		case "help":
			fmt.Println(usage)
			os.Exit(0)
		case "version":
			fmt.Println("Version: " + version)
			os.Exit(0)
		case "configure":
			c := shared.Configure()
			c.SaveToFile()
			os.Exit(0)
		case "start":
			fmt.Println("Starting...")
			os.Exit(0)
		case "node":
			fmt.Println("Starting as node...")
			host.Start()
			os.Exit(0)
		case "host":
			fmt.Println("Starting as host...")
			host.Start()
			os.Exit(0)
		default:
			fmt.Println("Unknown argument: " + arg)
			fmt.Println(usage)
			os.Exit(0)
		}
	}
}
