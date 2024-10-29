package mock

import (
	"reflect"
	"strings"
)

func findValue(m map[string]any, key string) any {
	keys := strings.Split(key, ".")

	current := m
	for _, v := range keys {
		if current[v] != nil && reflect.TypeOf(current[v]).Kind() == reflect.Map {
			current, _ = current[v].(map[string]any)
		} else {
			return current[v]
		}
	}

	return ""
}

func setValue(dest map[string]any, key string, value any) {
	keys := strings.Split(key, ".")

	current := dest

	for i := range keys {
		if i == len(keys)-1 {
			current[keys[i]] = value

			return
		}

		if current[keys[i]] == nil || reflect.TypeOf(current[keys[i]]).Kind() != reflect.Map {
			current[keys[i]] = map[string]any{}
		}

		current, _ = current[keys[i]].(map[string]any)
	}
}
