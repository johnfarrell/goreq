package goreq

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	tagName = "goreq"

	inKey       = "in"
	inValHeader = "header"
	inValQuery  = "query"

	labelKey = "label"
	typeKey  = "type"
)

type structTags struct {
	In    string
	Label string
	Type  string
}

// parseTags finds all the tag components and parses
// them into an easier to consume format
func parseTags(val string) structTags {
	values := strings.Split(val, ",")

	tags := structTags{}

	for _, value := range values {
		parts := strings.Split(value, "=")
		if len(parts) != 2 {
			continue
		}
		if len(parts[1]) == 0 {
			continue
		}

		switch parts[0] {
		case inKey:
			tags.In = parts[1]
		case labelKey:
			tags.Label = parts[1]
		case typeKey:
			tags.Type = parts[1]
		}
	}

	// Default to query parameter search if no valid inKey was
	// specified
	if tags.In == "" {
		tags.In = inValQuery
	}

	return tags
}

func parseParameters[T any](r *http.Request) (*T, error) {
	var ret T

	// Used for finding struct tag definitions
	ty := reflect.TypeOf(ret)

	// Make sure we only try to parse into a struct-type. We are
	// unable to enforce this via the generic type so need to do it here.
	if ty.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got a %s", ty)
	}

	// Used for assigning return values
	va := reflect.ValueOf(&ret)

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)

		tagVal := field.Tag.Get(tagName)

		// Don't process any skipped tags
		if tagVal == "-" {
			continue
		}
		tag := parseTags(tagVal)

		key := tag.Label
		if key == "" {
			key = field.Name
		}

		val := ""
		switch tag.In {
		case inValHeader:
			val = r.Header.Get(key)
		case inValQuery:
			val = r.URL.Query().Get(strings.ToLower(key))
		}

		// If the value was not found, return an error
		if val == "" {
			return nil, fmt.Errorf("field %s has no value for %s", field.Name, key)
		}

		va.Elem().FieldByName(field.Name).SetString(val)
	}

	return &ret, nil
}
