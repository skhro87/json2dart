package lib

import (
	"github.com/jeffail/gabs"
	"strings"
	"github.com/jinzhu/inflection"
	"fmt"
	"reflect"
)

func propTypeDart(propName string, prop interface{}) (string, error) {
	propType := ""
	switch prop.(type) {
	case string:
		propType = "String"
	case float64:
		propType = "num"
	case int:
		propType = "num"
	case map[string]interface{}, []interface{}:
		propType = childObjectClassNameFromPropName(propName)
	default:
		return "", fmt.Errorf("unsupported prop type : %v", reflect.TypeOf(prop))
	}
	return propType, nil
}

func linesFields(properties map[string]*gabs.Container) ([]string, error) {
	var lines []string
	for propName, prop := range properties {
		propType, err := propTypeDart(propName, prop.Data())
		if err != nil {
			return []string{}, fmt.Errorf("err getting prop type of %v : %v", prop.Data(), err.Error())
		}

		var line string
		switch prop.Data().(type) {
		case []interface{}:
			count, err := prop.ArrayCount()
			if err != nil {
				return []string{}, fmt.Errorf("err getting count from array : %v", err.Error())
			}

			if count > 0 {
				containerFirstInArray, err := prop.ArrayElement(0)
				if err != nil {
					return []string{}, fmt.Errorf("err getting first child of array : %v", err.Error())
				}

				switch containerFirstInArray.Data().(type) {
				case string, int, float64:
					propTypeFirstInArray, err := propTypeDart(propName, containerFirstInArray.Data())
					if err != nil {
						return []string{}, fmt.Errorf("err getting prop type : %v", err.Error())
					}
					line = fmt.Sprintf("final List<%v> %v;", propTypeFirstInArray, cleanPropName(propName))
				case interface{}, []interface{}:
					line = fmt.Sprintf("final List<%v> %v;", propType, cleanPropName(propName))
				}
			}
		default:
			line = fmt.Sprintf("final %v %v;", propType, cleanPropName(propName))
		}

		lines = append(lines, line)
	}

	return lines, nil
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

func linesFromJson(properties map[string]*gabs.Container) ([]string, error) {
	var lines []string
	for propNameRaw := range properties {
		propName := cleanPropName(propNameRaw)
		line := fmt.Sprintf("%v: json['%v'],", propName, propNameRaw)
		lines = append(lines, line)
	}

	return lines, nil
}

func childObjectClassNameFromPropName(propName string) string {
	propNamePartsRaw := strings.Split(cleanPropName(propName), "_")
	var propNameParts []string
	for _, part := range propNamePartsRaw {
		propName := inflection.Singular(strings.Title(strings.ToLower(part)))
		propNameParts = append(propNameParts, propName)
	}
	return strings.Join(propNameParts, "")
}

func cleanPropName(name string) string {
	return strings.Replace(name, " ", "_", -1)
}
