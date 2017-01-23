package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"text/template"
)

var templates = template.Must(template.ParseGlob("arduino/template_home/*.ino"))

type BoxTemplate struct {
	BoxID   string
	Sensors []SensorTemplate
}

type SensorTemplate struct {
	Define string
}

func generateSketchFromTemplate(box Box) (string, error) {

	// Define sensors
	var sensors []SensorTemplate
	var customSensorID = 1
	for _, v := range box.Sensors {
		switch v.Title {
		case "Temperatur":
			sensors = append(sensors, SensorTemplate{Define: "#define TEMPSENSOR_ID " + v.ID})
		case "rel. Luftfeuchte":
			sensors = append(sensors, SensorTemplate{Define: "#define HUMISENSOR_ID " + v.ID})
		case "Luftdruck":
			sensors = append(sensors, SensorTemplate{Define: "#define PRESSURESENSOR_ID " + v.ID})
		case "Lautstärke":
			sensors = append(sensors, SensorTemplate{Define: "#define NOISESENSOR_ID " + v.ID})
		case "Helligkeit":
			sensors = append(sensors, SensorTemplate{Define: "#define LIGHTSENSOR_ID " + v.ID})
		case "Beleuchtungsstärke":
			sensors = append(sensors, SensorTemplate{Define: "#define LUXSENSOR_ID " + v.ID})
		case "UV-Intensität":
			sensors = append(sensors, SensorTemplate{Define: "#define UVSENSOR_ID " + v.ID})
		default:
			sensors = append(sensors, SensorTemplate{Define: "#define SENSOR" + strconv.Itoa(customSensorID) + "_ID " + v.ID})
			customSensorID++
		}
	}

	boxTemplate := BoxTemplate{BoxID: box.ID, Sensors: sensors}

	var strBuffer bytes.Buffer
	err := templates.ExecuteTemplate(&strBuffer, "template_"+box.Model+".ino", boxTemplate)
	if err != nil {
		return "", err
	}

	// var boxID string
	str := strBuffer.String()
	fmt.Println(str)
	// for k, v := range payload.(map[string]interface{}) {
	// 	if k == "box" {
	// 		for k1, v1 := range v.(map[string]interface{}) {
	// 			fmt.Println(k1)
	// 			switch t := v1.(type) {
	// 			case string:
	// 				boxID = t
	// 			case []interface{}:
	// 				for k2, v2 := range v1.([]interface{}) {
	// 					fmt.Println(k2)
	// 					fmt.Println(v2)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	writeFileErr := ioutil.WriteFile("./"+box.ID+".ino", []byte(str), 0644)
	if writeFileErr != nil {
		log.Fatalln("[Templates - GenerateSketchFromTemplate] Writing INO file failed!")
		return "", writeFileErr
	}
	return "./" + box.ID + ".ino", nil
}
