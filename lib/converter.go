package lib

import (
	"github.com/jeffail/gabs"
	"fmt"
	"io/ioutil"
	"strings"
)

func Json2DartFile(json, rootClassName, fileLocation string) (error) {
	res, err := Json2Dart(json, rootClassName)
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

func Json2Dart(json, rootClassName string) (string, error) {
	if rootClassName == "" {
		return "", fmt.Errorf("root class name cannot be blank")
	}

	rootClassName = strings.Title(strings.ToLower(rootClassName))

	jsonParsed, err := gabs.ParseJSON([]byte(json))
	if err != nil {
		return "", fmt.Errorf("err parsing input json : %v", err.Error())
	}

	classes, err := containerToClasses(jsonParsed, rootClassName, []string{})
	if err != nil {
		return "", fmt.Errorf("err converting container to classes : %v", err.Error())
	}

	return strings.Join(classes, "\n\n"), nil
}

func containerToClasses(c *gabs.Container, className string, classes []string) ([]string, error) {
	nestedClasses, err := buildNestedClasses(c, classes)
	if err != nil {
		return []string{}, fmt.Errorf("err getting nested classes : %v", err.Error())
	}

	classes = append(classes, nestedClasses...)

	classCurrent, err := buildCurrentClass(c, className)
	if err != nil {
		return []string{}, fmt.Errorf("err converting current class : %v", err.Error())
	}

	classes = append(classes, classCurrent)

	return classes, nil
}

func buildCurrentClass(c *gabs.Container, className string) (string, error) {
	props, err := c.ChildrenMap()
	if err != nil {
		return "", fmt.Errorf("err parsing input json children : %v", err.Error())
	}

	linesFields, err := linesFields(props)
	if err != nil {
		return "", fmt.Errorf("err building props list : %v", err.Error())
	}

	constructor, err := constructor(className, props)
	if err != nil {
		return "", fmt.Errorf("err building constructor")
	}

	linesFromJson, err := linesFromJson(props)
	if err != nil {
		return "", fmt.Errorf("err building props assigment : %v", err.Error())
	}

	return buildClassFromParts(className, linesFields, constructor, linesFromJson), nil
}

func buildClassFromParts(className string, linesFields []string, constructor string, linesFromJson []string) string {
	res := fmt.Sprintf("class %v {", className)
	for _, l := range linesFields {
		res = fmt.Sprintf("%v\n\t%v", res, l)
	}
	res = fmt.Sprintf("%v\n\n\t%v", res, constructor)
	linesFactory := fmt.Sprintf("%v.fromJson(Map<String, dynamic> json) {\n\t\treturn new %v(", className, className)
	res = fmt.Sprintf("%v\n\n\t%v", res, linesFactory)
	for _, l := range linesFromJson {
		res = fmt.Sprintf("%v\n\t\t\t%v", res, l)
	}
	res = fmt.Sprintf("%v\n\t\t);\n\t}\n}", res)
	return res
}

func buildNestedClasses(c *gabs.Container, classes []string) ([]string, error) {
	props, err := c.ChildrenMap()
	if err != nil {
		return []string{}, fmt.Errorf("err parsing input json children : %v", err.Error())
	}

	var nestedClasses []string

	for propName, propContainer := range props {
		switch propContainer.Data().(type) {
		case map[string]interface{}:
			propNestedClasses, err := containerToClasses(propContainer, childObjectClassNameFromPropName(propName), classes)
			if err != nil {
				return []string{}, fmt.Errorf("err converting container to classes for prop %v : %v", propName, err.Error())
			}
			nestedClasses = append(nestedClasses, propNestedClasses...)
		case []interface{}:
			count, err := propContainer.ArrayCount()
			if err != nil {
				return []string{}, fmt.Errorf("err getting count from array : %v", err.Error())
			}

			if count > 0 {
				containerFirstInArray, err := propContainer.ArrayElement(0)
				if err != nil {
					return []string{}, fmt.Errorf("err getting first child of array : %v", err.Error())
				}

				propNestedClasses, err := containerToClasses(containerFirstInArray, childObjectClassNameFromPropName(propName), classes)
				if err != nil {
					return []string{}, fmt.Errorf("err converting container to classes for prop %v : %v", propName, err.Error())
				}
				nestedClasses = append(nestedClasses, propNestedClasses...)
			}
		}
	}

	return nestedClasses, nil
}
