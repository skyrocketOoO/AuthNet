package utils

import (
	"reflect"
	"strings"

	"github.com/skyrocketOoO/AuthNet/domain"
)

func ValidateRel(rel domain.Edge) error {
	if rel.ObjNs == "" || rel.ObjName == "" || rel.ObjRel == "" ||
		rel.SbjNs == "" || rel.SbjName == "" {
		return domain.ErrBodyAttribute{}
	}
	return ValidateReservedWord(rel)
}

func ValidateVertex(vertex domain.Vertex, isSubject bool) error {
	if vertex.Ns == "" || vertex.Name == "" {
		return domain.ErrBodyAttribute{}
	}
	if !isSubject && vertex.Rel == "" {
		return domain.ErrBodyAttribute{}
	}
	return ValidateReservedWord(vertex)
}

func ValidateReservedWord(st interface{}) error {
	value := reflect.ValueOf(st)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		str := field.Interface().(string)
		if strings.Contains(str, "%") {
			return domain.ErrBodyAttribute{}
		}

	}
	return nil
}
