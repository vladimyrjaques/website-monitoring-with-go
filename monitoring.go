package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitorings = 2
const delay = 5

func main() {
	showIntroduction()

	fmt.Printf("\n")

	for {
		sites := readWesitesFile()

		showMenu()

		command := selectCommand()

		switch command {
		case 1:
			startMonitoring(sites)
		case 2:
			showTodayLogs()
		case 3:
			addSitesToVerification()
		case 0:
			fmt.Println("Exiting program...")
			os.Exit(0)
		default:
			fmt.Println("Invalid command...")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	//Variable declaration types
	var name string = "Vladimyr"
	var lastName = "Jaques"
	version := 1.1

	//Verify variable type
	// fmt.Println("The version variable type is: ", reflect.TypeOf(version))

	fmt.Println("Hello Mr.", name, lastName)
	fmt.Println("This program is in version: ", version)
}

func showMenu() {
	fmt.Println("\n1 - Start monitoring")
	fmt.Println("2 - Show today's logs")
	fmt.Println("3 - Add sites to verification")
	fmt.Println("0 - Exit Program")
}

func selectCommand() int {
	var command int
	fmt.Scan(&command)

	// fmt.Println("\n\nThe address my variable command is: ", &command)
	fmt.Println("Chosen command: ", command)

	return command
}

func startMonitoring(sites []string) {
	fmt.Println("Monitoring...")

	for i := 0; i < monitorings; i++ {
		for _, site := range sites {
			monitorWebsite(site)
		}

		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func monitorWebsite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("An error has occurred: ", err)

		return
	}

	if resp.StatusCode == 200 {
		fmt.Println("\nSite: ", site, "successfully loaded")

		return
	}

	fmt.Println("Site: ", site, "it has problems")
	registerLogs(site, false, "Not loaded")
}

func readWesitesFile() []string {
	var sites []string
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("An error has occurred in file read")
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

func registerLogs(site string, status bool, message string) {
	file, err := os.OpenFile("./logs/log-"+time.Now().Format("2006-01-02")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	file.WriteString("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + site + ": status " + strconv.FormatBool(status) + ", " + message + "\n")
	file.Close()
}

func showTodayLogs() {
	fmt.Println("\nShowing today's logs...")
	file, err := ioutil.ReadFile("./logs/log-" + time.Now().Format("2006-01-02") + ".txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\n", string(file))
}

func addSitesToVerification() {
	var sites string

	fmt.Println("\nEnter sites separated by comma: ")
	fmt.Scan(&sites)

	splitSites := strings.Split(sites, ",")
	file, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	var invalidSites []int

	for key, site := range splitSites {
		if isUrl(site) {
			file.WriteString("\n" + site)

			continue
		}

		invalidSites = append(invalidSites, key)
	}

	file.Close()

	fmt.Println("\nInvalid sites not added: ")

	for _, siteKey := range invalidSites {
		fmt.Println(splitSites[siteKey])
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
