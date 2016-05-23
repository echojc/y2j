package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) > 2 {
		fmt.Println("Usage: y2j [file]")
		fmt.Println("If file is omitted, then read from stdin.")
		os.Exit(2)
	}

	input := os.Stdin
	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		input = file
	}
	defer input.Close()

	y, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(y, &data)
	if err != nil {
		log.Fatal(err)
	}

	normalized, err := NormalizeObject(data)
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.Marshal(normalized)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Compact(&out, j)
	out.WriteTo(os.Stdout)
	fmt.Println()
}

func NormalizeObject(data map[interface{}]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for key, value := range data {
		stringKey, success := key.(string)
		if !success {
			return nil, errors.New(fmt.Sprintf("Key was not a string: %v\n", key))
		}

		if mapValue, success := value.(map[interface{}]interface{}); success {
			normalized, err := NormalizeObject(mapValue)
			if err != nil {
				return nil, err
			}
			out[stringKey] = normalized
		} else if arrayValue, success := value.([]interface{}); success {
			normalized, err := NormalizeArray(arrayValue)
			if err != nil {
				return nil, err
			}
			out[stringKey] = normalized
		} else {
			out[stringKey] = value
		}
	}
	return out, nil
}

func NormalizeArray(data []interface{}) ([]interface{}, error) {
	out := make([]interface{}, 0, len(data))
	for _, value := range data {
		if mapValue, success := value.(map[interface{}]interface{}); success {
			normalized, err := NormalizeObject(mapValue)
			if err != nil {
				return nil, err
			}
			out = append(out, normalized)
		} else if arrayValue, success := value.([]interface{}); success {
			normalized, err := NormalizeArray(arrayValue)
			if err != nil {
				return nil, err
			}
			out = append(out, normalized)
		} else {
			out = append(out, value)
		}
	}
	return out, nil
}
