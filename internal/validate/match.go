package validate

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tmavrin/mock-http-server/internal/config"
)

func bodyMatching(httpReq any, req *config.Request) error {
	var expected any

	if err := json.Unmarshal([]byte(req.Match), &expected); err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(httpReq)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, httpReq) {
		return fmt.Errorf(
			"request does not match expected: %s , actual: %s",
			req.Match,
			string(jsonBytes),
		)
	}

	return nil
}
