package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 3
const delay = 5

func main() {
	showIntroduction()

	for {
		showMenu()

		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Showing logs...")
			printLogs()
		case 0:
			fmt.Println("Exiting....")
			os.Exit(0)
		default:
			fmt.Println("I don't know this command")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	name := "Daniel"
	version := 1.1
	fmt.Println("Hello, mr.", name)
	fmt.Println("This program is in the version:", version)
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	fmt.Println(command)
	return command
}

func showMenu() {
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Exit")
}

func startMonitoring() {
	fmt.Println("Starting monitoring...")
	sites := readSitesFromFile()

	for i := 0; i < monitoring; i++ {
		fmt.Println("Testando sites", i+1, " de 5")
		for _, site := range sites {
			testSite(site)
		}
		fmt.Println("-----------------------------")
		time.Sleep(delay * time.Second)
	}
}

func testSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro ao acessar o site", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status code", resp.StatusCode)
		registerLog(site, false)
	}
}

func readSitesFromFile() []string {
	var sites []string
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro ", err)
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return sites
}

func registerLog(site string, online bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo", err)
	}

	now := time.Now().Format("02/01/2006 15:04:05")
	file.WriteString(now + " " + site + " - Online: " + strconv.FormatBool(online) + "\n")
	file.Close()
}

func printLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo")
	}

	fmt.Println(string(file))
}
