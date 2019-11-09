package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	for {
		showMenu()

		command := readCommand()

		switch command {
		case 1:
			fmt.Println("Iniciando monitoramento...")
			initMonitor()
		case 2:
			fmt.Println("Exibindo LOGs...")
			logsReader()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func showMenu() {

	fmt.Println("1- Iniciar monitoramento")
	fmt.Println("2- Exibir LOGs")
	fmt.Println("0- Sair do programa")

}

func readCommand() int {
	var command int

	fmt.Scan(&command)

	return command
}

func initMonitor() {

	sites := readSites()

	for _, site := range sites {
		testSites(site)
		fmt.Println("")
	}
}

func readSites() []string {

	var sites []string

	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro :", err)
		os.Exit(-1)
	}

	buffer := bufio.NewReader(file)
	for {
		line, err := buffer.ReadString('\n')
		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return sites
}

func testSites(site string) {
	response, err := http.Get(site)

	if err != nil {
		fmt.Println("Erro: ", err)
	}

	if response.StatusCode == 200 {
		fmt.Println("Site: ", site, "carregado com sucesso!")
		logsRegister(site, true)
	} else {
		fmt.Println("Site: ", site, " com problemas, Status Code: ", response.StatusCode)
		logsRegister(site, false)
	}
}

func logsRegister(site string, status bool) {

	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	var statusString string
	if status {
		statusString = "online"
	} else {
		statusString = "offline"
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " Status: " + statusString + "\n")
	file.Close()
}

func logsReader() {
	file, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(string(file))
}
