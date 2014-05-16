package main

import (
        "log"
        "fmt"
        "flag"
        "net/http"
        "crypto/tls"
        "io/ioutil"
        "strings"
        _ "neuroinformatics.harvard.edu/survana"
       )

const (
        DEFAULT_STORE_URL = "https://localhost:4443/store/"
        DOWNLOAD_ENDPOINT = "download"
      )

var (
        store_url string
        study_id  string
    )

func main() {
    ParseArguments()

    log.Println("Requesting all responses for study", study_id)
    data, err := download(store_url + DOWNLOAD_ENDPOINT + "?" + study_id)
    if err != nil {
        log.Fatalln(err)
    }

    log.Println("data=", string(data))
}

func ParseArguments() {

	//command line arguments
	flag.StringVar(&store_url, "store", DEFAULT_STORE_URL, "Survana Store URL")
	flag.StringVar(&study_id, "study", "", "Survana Study ID")
	flag.Parse()

    //make sure store_url is terminated by a slash
    if (!strings.HasSuffix(store_url, "/")) {
        store_url += "/"
    }
}

func download(url string) (result []byte, err error) {

    tr := &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify : true,
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
