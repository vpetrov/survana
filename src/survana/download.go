package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	DEFAULT_STORE_URL = "https://localhost:4443/store/"
	DOWNLOAD_ENDPOINT = "download"
)

var (
	store_url   string
	study_id    string
	output_file string
)

func main() {
	var (
		err error
	)

	ParseArguments()

	log.Println("Requesting all responses for study", study_id, "from ", store_url)
	data, err := download(store_url + DOWNLOAD_ENDPOINT + "/" + study_id)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(output_file, data, 0400)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("All responses saved to", output_file)
}

func ParseArguments() {

	//command line arguments
	flag.StringVar(&store_url, "store", DEFAULT_STORE_URL, "Survana Store URL")
	flag.StringVar(&study_id, "study", "", "Survana Study ID")
	flag.StringVar(&output_file, "output", "responses.json", "Path to output file")
	flag.Parse()

	if len(study_id) == 0 {
		log.Fatalln("Error: the -study flag is required. See --help.")
	}

	//make sure store_url is terminated by a slash
	if !strings.HasSuffix(store_url, "/") {
		store_url += "/"
	}

	if _, err := os.Stat(output_file); err == nil {
		log.Fatalln("Error:", output_file, ": file exists")
	}
}

func download(url string) (result []byte, err error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	response, err := client.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	//only return a non-nil result when http status is 200, otherwise use the body as the message string
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("download error: %s", string(result))
		result = nil
		return
	}

	return
}
