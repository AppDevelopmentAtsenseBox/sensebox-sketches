package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
)

var templates = template.Must(template.ParseGlob("arduino/template_home/*.ino"))

func generateSketchFromTemplate(networkType string, payload interface{}) (string, error) {
	var strBuffer bytes.Buffer
	err := templates.ExecuteTemplate(&strBuffer, "template_home_"+networkType+".ino", payload)
	if err != nil {
		return "", err
	}

	var boxID string
	str := strBuffer.String()
	for k, v := range payload.(map[string]interface{}) {
		if k == "box" {
			for k1, v1 := range v.(map[string]interface{}) {
				fmt.Println(k1)
				switch t := v1.(type) {
				case string:
					boxID = t
				case []interface{}:
					for k2, v2 := range v1.([]interface{}) {
						fmt.Println(k2)
						fmt.Println(v2)
					}
				}
			}
		}
	}
	writeFileErr := ioutil.WriteFile("./"+boxID+".ino", []byte(str), 0644)
	if writeFileErr != nil {
		log.Fatalln("[Templates - GenerateSketchFromTemplate] Writing INO file failed!")
		return "", writeFileErr
	}
	return "./" + boxID + ".ino", nil
}
