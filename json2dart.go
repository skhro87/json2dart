package main

import (
	"os"
	"log"
	"io/ioutil"
	"github.com/skhro87/json2dart/lib"
	"strings"
	"fmt"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "json2dart"
	app.Version = "0.1.0"
	app.Usage = "convert json to dart code"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Usage: "specify input file (if not using pipe)",
		},
		cli.StringFlag{
			Name:  "out",
			Value: "./model_gen",
			Usage: "specify output file(s) folder",
		},
		cli.StringFlag{
			Name:  "class",
			Value: "Root",
			Usage: "name of root class",
		},
		cli.BoolFlag{
			Name: "split",
			Usage: "split generated classes into multiple files",
		},
		cli.BoolFlag{
			Name: "no-files",
			Usage: "disable output to file(s)",
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func run(c *cli.Context) error {
	inputFile := c.String("in")
	outputFolder := c.String("out")
	rootClassName := c.String("class")
	splitClasses := c.Bool("split")
	noFileOutput := c.Bool("no-files")

	input, err := readInput(inputFile)
	if err != nil {
		return fmt.Errorf("err reading input : %v", err.Error())
	}

	classes, err := lib.Json2Dart(input, rootClassName)
	if err != nil {
		log.Fatalf("err converting : %v", err.Error())
	}

	if !noFileOutput && outputFolder != "" {
		err = writeOutput(classes, outputFolder, splitClasses)
		if err != nil {
			return fmt.Errorf("err writing classes to file : %v", err.Error())
		}
	}

	printOutput(classes)

	return nil
}

func readInput(inputFile string) (string, error) {
	if inputFile != "" {
		input, err := inputFromFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("err reading input from file %v : %v", inputFile, err.Error())
		}
		return input, nil
	}

	input := inputFromPipe()
	if input == "" {
		return "", fmt.Errorf("no input file specified, and no input via pipe, please use one of both options")
	}

	return input, nil
}

func inputFromFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	return string(bytes), err
}

func inputFromPipe() (string) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return ""
	}

	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("err reading from stdin : %v", err.Error())
	}

	return strings.TrimSpace(string(inputBytes))
}

func printOutput(classes []lib.ClassDef) {
	for _, class := range classes {
		fmt.Printf("%v\n\n", class.Code)
	}
}

func writeOutput(classes []lib.ClassDef, folder string, split bool) error {
	if folder == "" {
		return fmt.Errorf("blank folder name")
	}

	_, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(folder, 0644)
			if err != nil {
				return fmt.Errorf("err creating folder %v : %v", folder, err.Error())
			}
		} else {
			return fmt.Errorf("err checking folder %v : %v", folder, err.Error())
		}
	}

	if split {
		return outputToMultipleFiles(classes, folder)
	}
	return outputToSingleFile(classes, folder)
}

func outputToSingleFile(classes []lib.ClassDef, folder string) error {
	file := fmt.Sprintf("%v/%v.dart", folder, "model")

	output := ""
	for _, class := range classes {
		output = fmt.Sprintf("%v\n\n%v", output, class.Code)
	}

	err := ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("err writing file to %v : %v", file, err.Error())
	}

	return nil
}

func outputToMultipleFiles(classes []lib.ClassDef, folder string) error {
	for _, class := range classes {
		file := fmt.Sprintf("%v/%v.dart", folder, strings.ToLower(class.ClassName))

		err := ioutil.WriteFile(file, []byte(class.Code), 0644)
		if err != nil {
			return fmt.Errorf("err writing file to %v : %v", file, err.Error())
		}
	}

	return nil
}
