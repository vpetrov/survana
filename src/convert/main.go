package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"neuroinformatics.harvard.edu/survana"
	"os"
	_ "reflect"
)

func main() {

	var (
		filename  string
		err       error
		file      *os.File
		old_forms []OldForm
		new_forms []survana.Form
	)

	flag.StringVar(&filename, "f", "", "Path to filename containing an array of JSON forms")
	flag.Parse()

	if file, err = os.Open(filename); err != nil {
		log.Fatalln(err)
	}

	if old_forms, err = read_old_forms(file); err != nil {
		print_json_error(err)
		log.Fatalln(err)
	}

	log.Println("Found", len(old_forms), "old forms.")

	if new_forms, err = convert_forms(old_forms); err != nil {
		log.Fatalln(err)
	}

	log.Println("Converted into", len(new_forms), "new forms.")
}

func read_old_forms(file *os.File) (old_forms []OldForm, err error) {
	decoder := json.NewDecoder(file)
	var i int

	for {
		i++
		var old_form OldForm
		err = decoder.Decode(&old_form)
		if err != nil {
			if err == io.EOF {
				return old_forms, nil
			}

			log.Println("line", i)
			return
		}

		old_forms = append(old_forms, old_form)
	}

	return
}

func convert_forms(old_forms []OldForm) (forms []survana.Form, err error) {

	for i, old_form := range old_forms {
		log.Printf("%v: %#v", i, old_form)
		_ = survana.Form{}
		break
	}

	return
}

func print_json_error(err error) {
	switch err := err.(type) {
	default:
		log.Println(err)
	case *json.InvalidUTF8Error:
		log.Printf("%s: '%s'\n", err, err.S)
	case *json.InvalidUnmarshalError:
		log.Printf("%s: %s\n", err.Type, err)
	case *json.MarshalerError:
		log.Printf("%s: %s\n", err.Type, err.Err)
	case *json.SyntaxError:
		log.Printf("syntax error: offset %v: %s\n", err.Offset, err, err)
	case *json.UnmarshalFieldError:
		log.Printf("key '%s', type '%s': %s\n", err.Key, err.Type, err)
	case *json.UnmarshalTypeError:
		log.Printf("value '%s', type '%s': %s\n", err.Value, err.Type, err)
	case *json.UnsupportedTypeError:
		log.Printf("type '%s': %s\n", err.Type, err)
	case *json.UnsupportedValueError:
		log.Printf("value '%s', string '%s': %s\n", err.Value, err.Str, err)
	}
}
