//sends old recovered appcache-item elements to survana stores
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type KeyInfo struct {
	Id   string `json:"id"`
	PEM  string `json:"pem"`
	Bits string `json:"bits"`
}

type Payload struct {
	Password string `json:"password"`
	Data     string `json:"data"`
}

type Response struct {
	Data struct {
		Key     KeyInfo `json:"key"`
		Payload Payload `json:"payload"`
	} `json:"data"`
	Url *url.URL `json:"url"`
}

type ResponseString struct {
	Data string
	Url  string
}

func main() {

	var (
		filepath string
		filedata []byte
		err      error
	)

	flag.StringVar(&filepath, "f", "", "Path to recovery file")
	flag.Parse()

	if filedata, err = ioutil.ReadFile(filepath); err != nil {
		log.Fatalf("%s: %s\n", filepath, err)
	}

	lines := strings.Split(string(filedata), "\n")
	if len(lines) == 0 {
		log.Fatalf("%s: no lines found")
	}

	log.Println("Found?", len(lines), "lines in file.")

	type ParsingFailure struct {
		Index int
		Line  string
		Error error
	}

	var (
		failed    []ParsingFailure
		succeeded int
	)

	for i, line := range lines {
		response_id, response, err := parseLine(line)
		if err != nil {
			log.Printf("%s: ", response_id)
			log.Println(err)
			failed = append(failed, ParsingFailure{i, line, err})
			continue
		}

		//skip invalid lines
		if len(response_id) == 0 {
			continue
		}

		err = sendResponse(response_id, response)
		if err != nil {
			failed = append(failed, ParsingFailure{i, line, err})
		}

		log.Println(response_id, "OK")

		succeeded++
	}

	log.Println("Succeeded: ", succeeded)
	log.Println("Failed: ", len(failed), failed)
}

func parseLine(line string) (response_id string, response *Response, err error) {
	if !strings.HasPrefix(line, "appcache-item-") {
		return
	}

	id_data := strings.SplitN(line, ":", 2)

	if len(id_data) < 2 {
		err = errors.New("Invalid line")
		return
	}

	response_id = strings.TrimSpace(id_data[0])

	if len(response_id) == 0 {
		err = errors.New("Invalid response ID")
		return
	}

	json_data := []byte(id_data[1])

	if len(json_data) == 0 {
		err = errors.New("Invalid data")
		return
	}

	response_string := &ResponseString{}

	if err = json.Unmarshal(json_data, response_string); err != nil {
		return
	}

	response = &Response{}
	if err = json.Unmarshal([]byte(response_string.Data), &response.Data); err != nil {
		return
	}

	if response.Url, err = url.Parse(response_string.Url); err != nil {
		return
	}

	return
}

func sendResponse(response_id string, response *Response) (err error) {

	var data []byte

	if data, err = json.Marshal(response.Data); err != nil {
		return
	}

	values := response.Url.Query()
	values.Set("callback", "send")
	values.Set("id", response_id)
	values.Set("data", string(data))

	data_url := response.Url.String() + "?" + values.Encode()
	log.Println("Sending Response", response_id, "to", data_url)
	var resp_stream *http.Response
	if resp_stream, err = http.Get(data_url); err != nil {
		return
	}

	defer func() { resp_stream.Body.Close() }()

	var resp []byte
	if resp, err = ioutil.ReadAll(resp_stream.Body); err != nil {
		return
	}

	log.Println("response=", string(resp))

	return nil
}
