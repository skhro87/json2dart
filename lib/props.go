package lib

import (
	"github.com/jeffail/gabs"
	"strings"
	"github.com/jinzhu/inflection"
	"fmt"
	"reflect"
)

func typeOfProp(propName string, prop interface{}) (string, error) {
	propType := ""
	switch prop.(type) {
	case string:
		propType = "String"
	case float64:
		propType = "double"
	case int:
		propType = "int"
	case map[string]interface{}:
		propType = childObjectClassNameFromPropName(propName)
	default:
		return "", fmt.Errorf("unsupported prop type : %v", reflect.TypeOf(prop))
	}
	return propType, nil
}

func propsList(properties map[string]*gabs.Container) (string, error) {
	out := ``
	for propName, prop := range properties {
		propType, err := typeOfProp(propName, prop.Data())
		if err != nil {
			return "", fmt.Errorf("err getting prop type of %v : %v", prop.Data(), err.Error())
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
