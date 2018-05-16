package main

import (
	"os"
	"log"
	"io/ioutil"
	"github.com/skhro87/json2dart/lib"
	"strings"
	"fmt"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("err getting stdin stat : %v", err.Error())
	}
	if fi.Mode() & os.ModeNamedPipe == 0 {
		log.Printf("no input :/")
		return
	}

	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("err reading from stdin : %v", err.Error())
	}

	if strings.TrimSpace(string(inputBytes)) == "" {
		log.Fatalf("empty input :/")
	}

	log.Printf("converting...")

	res, err := lib.Json2Dart(string(inputBytes), "Root")
	if err != nil {
		log.Fatalf("err converting : %v", err.Error())
	}

	log.Printf("done!")
	log.Printf("\n%v", res)
}

func toFile(input, rootClassName string, fileLocation string) error {
	res, err := lib.Json2Dart(input, rootClassName)
	if err != nil {
		return fmt.Errorf("err creating dart code from json : %v", err.Error())
	}

	filename := fmt.Sprintf("%v.dart", rootClassName)

	if fileLocation == "" {
		fileLocation = fmt.Sprintf("./out/%v", filename)
	} else {
		fileLocation = fmt.Sprintf("%v/%v", fileLocation, filename)
	}

	err = ioutil.WriteFile(fileLocation, []byte(res), 0644)
	if err != nil {
		return fmt.Errorf("err writing file to %v : %v", fileLocation, err.Error())
	}

	return nil
}
