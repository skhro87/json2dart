package lib

import (
	"github.com/jeffail/gabs"
	"fmt"
	"strings"
)

func Json2Dart(json, rootClassName string) ([]ClassDef, error) {
	if rootClassName == "" {
		return []ClassDef{}, fmt.Errorf("root class name cannot be blank")
	}

	rootClassName = strings.Title(strings.ToLower(rootClassName))

	jsonParsed, err := gabs.ParseJSON([]byte(json))
	if err != nil {
		return []ClassDef{}, fmt.Errorf("err parsing input json : %v", err.Error())
	}

	classes, err := containerToClasses(jsonParsed, rootClassName, []ClassDef{})
	if err != nil {
		return []ClassDef{}, fmt.Errorf("err converting container to classes : %v", err.Error())
	}

	return classes, nil
}

func containerToClasses(c *gabs.Container, className string, classes []ClassDef) ([]ClassDef, error) {
	nestedClasses, err := buildNestedClasses(c, classes)
	if err != nil {
		return []ClassDef{}, fmt.Errorf("err getting nested classes : %v", err.Error())
	}

	classes = append(classes, nestedClasses...)

	classCurrent, err := buildCurrentClass(c, className)
	if err != nil {
		return []ClassDef{}, fmt.Errorf("err converting current class : %v", err.Error())
	}

	classes = append(classes, classCurrent)

	return classes, nil
}

func buildCurrentClass(c *gabs.Container, className string) (ClassDef, error) {
	props, err := c.ChildrenMap()
	if err != nil {
		return ClassDef{}, fmt.Errorf("err parsing input json children : %v", err.Error())
	}

	linesFields, err := linesFields(props)
	if err != nil {
		return ClassDef{}, fmt.Errorf("err building props list : %v", err.Error())
	}

	constructor, err := constructor(className, props)
	if err != nil {
		return ClassDef{}, fmt.Errorf("err building constructor")
	}

	linesFromJson, err := linesFromJson(props)
	if err != nil {
		return ClassDef{}, fmt.Errorf("err building props assigment : %v", err.Error())
	}

	code := buildClassCode(className, linesFields, constructor, linesFromJson)

	return ClassDef{
		Code: code,
		ClassName: className,
	}, nil
}

func buildClassCode(className string, linesFields []string, constructor string, linesFromJson []string) string {
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

func buildNestedClasses(c *gabs.Container, classes []ClassDef) ([]ClassDef, error) {
	props, err := c.ChildrenMap()
	if err != nil {
		return []ClassDef{}, fmt.Errorf("err parsing input json children : %v", err.Error())
	}

	var nestedClasses []ClassDef

	for propName, propContainer := range props {
		switch propContainer.Data().(type) {
		case map[string]interface{}:
			propNestedClasses, err := containerToClasses(propContainer, childObjectClassNameFromPropName(propName), classes)
			if err != nil {
				return []ClassDef{}, fmt.Errorf("err converting container to classes for prop %v : %v", propName, err.Error())
			}
			nestedClasses = append(nestedClasses, propNestedClasses...)
		case []interface{}:
			count, err := propContainer.ArrayCount()
			if err != nil {
				return []ClassDef{}, fmt.Errorf("err getting count from array : %v", err.Error())
			}

			if count > 0 {
				containerFirstInArray, err := propContainer.ArrayElement(0)
				if err != nil {
					return []ClassDef{}, fmt.Errorf("err getting first child of array : %v", err.Error())
				}

				propNestedClasses, err := containerToClasses(containerFirstInArray, childObjectClassNameFromPropName(propName), classes)
				if err != nil {
					return []ClassDef{}, fmt.Errorf("err converting container to classes for prop %v : %v", propName, err.Error())
				}
				nestedClasses = append(nestedClasses, propNestedClasses...)
			}
		}
	}

	return nestedClasses, nil
}
