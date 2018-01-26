package main

import (
	"encoding/xml"
	"log"
	"os"
)

var xmlHeader = `<?xml version="1.0"?>` + "\n"

func getFilePointer() *os.File {
	// Open a new file for writing only
	file, err := os.OpenFile(
		"1UserLogFile.xml",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	file.Write([]byte(xmlHeader))
	return file
}

func getXMLEventString(loggingObject interface{}) []byte {

	var xmlstring []byte
	if xmlstring, err := xml.MarshalIndent(loggingObject, "", "    "); err == nil {
		xmlstring = []byte(string(xmlstring))
		return xmlstring
	}
	return xmlstring

}

func logging(loggingObject interface{}, file *os.File) {
	xmlstring := getXMLEventString(loggingObject)
	_, err := file.Write(xmlstring)
	check(err)
}
