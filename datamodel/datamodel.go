package datamodel

import "reflect"

func IsStoredInDatabase(x reflect.StructField) bool {
	return x.Name != "Project" && x.Name != "App"
}
