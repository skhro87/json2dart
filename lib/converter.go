package lib

import (
	"github.com/Jeffail/gabs"
	"fmt"
	"reflect"
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

	children, err := jsonParsed.ChildrenMap()
	if err != nil {
		return "", fmt.Errorf("err parsing input json children : %v", err.Error())
	}

	return createClass(rootClassName, children)
}

func createClass(name string, properties map[string]*gabs.Container) (string, error) {
	if name == "" {
		name = "XXX"
	}

	propsHead, err := propsList(properties)
	if err != nil {
		return "", fmt.Errorf("err building props list : %v", err.Error())
	}

	constructor, err := constructor(name, properties)
	if err != nil {
		return "", fmt.Errorf("err building constructor")
	}

	propsAssignment, err := propsAssignment(properties)
	if err != nil {
		return "", fmt.Errorf("err building props assigment : %v", err.Error())
	}

	return fmt.Sprintf(`
		class %v {
			%v
			
			%v

			factory %v.fromJson(Map<String, dynamic> json) {
				return new %v(
					%v
				);
			}
		}
	`, name, propsHead, constructor, name, name, propsAssignment), nil
}

func propsList(properties map[string]*gabs.Container) (string, error) {
	out := ``
	for propName, v := range properties {
		propType, err := typeOfProp(v.Data())
		if err != nil {
			return "", fmt.Errorf("err getting prop type of %v : %v", v.Data(), err.Error())
		}

		out = fmt.Sprintf(`%v
			final %v %v;`, out, propType, strings.Replace(propName, " ", "_", -1))
	}

	return out, nil
}

func constructor(className string, properties map[string]*gabs.Container) (string, error) {
	out := fmt.Sprintf("%v({", className)
	for propName := range properties {
		out = fmt.Sprintf(`%vthis.%v,`, out, cleanPropName(propName))
	}
	out = strings.TrimRight(out, ",")
	out = fmt.Sprintf("%v})", out)

	return out, nil
}

func propsAssignment(properties map[string]*gabs.Container) (string, error) {
	out := ""
	for propNameRaw := range properties {
		propName := cleanPropName(propNameRaw)
		out = fmt.Sprintf(`%v
					%v: json['%v'],`, out, propName, propNameRaw)
	}

	return out, nil
}

func typeOfProp(prop interface{}) (string, error) {
	propType := ""
	switch prop.(type) {
	case string:
		propType = "String"
	case float64:
		propType = "double"
	case int:
		propType = "int"
	default:
		return "", fmt.Errorf("invalid head prop type : %v", reflect.TypeOf(prop))
	}
	return propType, nil
}

func cleanPropName(name string) string {
	return strings.Replace(name, " ", "_", -1)
}