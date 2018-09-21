package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Info ...
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// Swagger ...
type Swagger struct {
	Swagger     string                 `json:"swagger"`
	Info        Info                   `json:"info"`
	Host        string                 `json:"host"`
	Schemes     []string               `json:"schemes"`
	Consumes    []string               `json:"consumes"`
	Produces    []string               `json:"produces"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

func main() {
	var location string
	var version string
	var title string
	var host string
	var apiver string
	var consumes string
	var produces string
	var output string

	flag.StringVar(&location, "path", "./", "the directory path")
	flag.StringVar(&version, "swagger", "2.0", "swagger version")
	flag.StringVar(&title, "title", "Swagger gen", "Project title")
	flag.StringVar(&host, "host", "localhost", "API host")
	flag.StringVar(&apiver, "apiver", "1.0", "api version")
	flag.StringVar(&consumes, "consumes", "application/json", "encoding of requests")
	flag.StringVar(&produces, "produces", "application/json", "encoding of responses")
	flag.StringVar(&output, "output", "./", "output location")
	flag.Parse()

	if location == "" {
		location = "./"
	}

	var swaggers Swagger
	swaggers.Paths = map[string]interface{}{}
	swaggers.Definitions = map[string]interface{}{}
	swaggers.Schemes = []string{"http", "https"}

	if version == "" {
		version = "2.0"
	}
	swaggers.Swagger = version

	if title == "" {
		title = "Swagger"
	}

	if apiver == "" {
		apiver = "1.0"
	}

	swaggers.Info = Info{
		Title:   title,
		Version: apiver,
	}

	if host == "" {
		host = "localhost"
	}
	swaggers.Host = host

	if consumes == "" {
		consumes = "application/json"
	}
	swaggers.Consumes = []string{consumes}

	if produces == "" {
		consumes = "application/json"
	}
	swaggers.Produces = []string{produces}

	if output == "" {
		output = "./swagger.json"
	}

	files, err := ioutil.ReadDir(location)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		jsonFile, err := os.Open(location + "/" + f.Name())
		if err != nil {
			log.Fatalln(err)
			continue
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		var swagger Swagger

		err = json.Unmarshal(byteValue, &swagger)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		for key, value := range swagger.Paths {
			if _, ok := swaggers.Paths[key]; !ok {
				swaggers.Paths[key] = value
			}
		}

		for key, value := range swagger.Definitions {
			if _, ok := swaggers.Definitions[key]; !ok {
				swaggers.Definitions[key] = value
			}
		}

	}
	marshaled, err := json.MarshalIndent(&swaggers, "", "    ")
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	_, err = jsonFile.Write(marshaled)
	if err != nil {
		panic(err)
	}

	err = jsonFile.Close()
	if err != nil {
		panic(err)
	}

}
