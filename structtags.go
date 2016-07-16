package logfmt

import (
	"reflect"
	"strings"
)

func extractStructTag(t reflect.StructField) (string, string) {
	tags := strings.SplitN(t.Tag.Get("logfmt"), ",", 2)
	switch len(tags) {
	case 1:
		return tags[0], ""
	case 2:
		return tags[0], tags[1]
	default:
		return "", ""
	}
}
