package main

import (
	"fmt"
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

// Config is the server config populated from config files and/or environment variables
type Config struct {
	IP   string
	Port int
}

var config Config

func main() {
	err := gonfig.GetConf(configFile(), &config)
	if err != nil {
		log.Fatalln(err)
		os.Exit(100)
	}
	addr := config.IP + ":" + strconv.Itoa(config.Port)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/process", ParseForm)
	http.HandleFunc("/", ServeStatic)

	http.ListenAndServe(addr, nil)
	log.Println("Listening at " + addr + "...")

	//ServePage(addr)
	//ListenForm()
	fmt.Println("channelpack starting...")

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
