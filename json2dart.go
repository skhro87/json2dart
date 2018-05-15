package main

import (
	"os"
	"log"
	"io/ioutil"
	"github.com/skhro87/json2dart/lib"
	"strings"
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

	res, err := lib.Json2Dart(string(inputBytes), "XXX")
	if err != nil {
		log.Fatalf("err converting : %v", err.Error())
	}

	log.Printf("done!")
	log.Printf("\n%v", res)
}
