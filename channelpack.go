package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/tkanos/gonfig"
)

// config is the server config populated from config files and/or environment variables
type config struct {
	IP   string
	Port int
}

var myConfig config

// myPackWorker is the main worker for the whole process.  This tool was designed to be used by me so there is only one...
var myPackWorker *packWorker

func main() {
	err := gonfig.GetConf(configFile(), &myConfig)
	if err != nil {
		log.Fatalln(err)
		os.Exit(100)
	}
	addr := myConfig.IP + ":" + strconv.Itoa(myConfig.Port)

	log.Println("Parsing templates")
	initTemplates()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/thumb/", serveThumbnail)
	http.HandleFunc("/upload", parseUpload)
	http.HandleFunc("/process", parseProcess)
	http.HandleFunc("/remove", parseRemove)
	http.HandleFunc("/", serveStatic)

	log.Println("Starting pack worker")
	myPackWorker = newPackWorker()

	log.Println("Listening at " + addr + "...")
	http.ListenAndServe(addr, nil)
}

func configFile() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	filename := []string{"config.", env, ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	configFilePath := path.Join(filepath.Dir(dirname), "config", strings.Join(filename, ""))
	return configFilePath
}
