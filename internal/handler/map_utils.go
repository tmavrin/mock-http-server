package handler

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var indexRegex = regexp.MustCompile(`\[(.*?)\]`)

func findValue(m any, key string) any {
	keys := strings.Split(key, ".")

	for _, v := range keys {
		switch m := m.(type) {
		case map[string]any:
			return findValue(m[v], strings.Join(keys[1:], "."))
		case []any:
			search := indexRegex.FindStringSubmatch(v)
			if len(search) == 0 {
				return nil
			}

			index, err := strconv.Atoi(search[1])
			if err != nil {
				return nil
			}

			nextKey := key
			if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
				nextKey = strings.Join(keys[1:], ".")
			}

			if found := findValue(m[index], nextKey); found != nil {
				return found
			}
		case []map[string]any:
			search := indexRegex.FindStringSubmatch(v)
			if len(search) == 0 {
				return nil
			}

			index, err := strconv.Atoi(search[1])
			if err != nil {
				return nil
			}

			nextKey := key
			if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
				nextKey = strings.Join(keys[1:], ".")
			}

			if found := findValue(m[index], nextKey); found != nil {
				return found
			}
		default:
			return m
		}
	}

	return nil
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
