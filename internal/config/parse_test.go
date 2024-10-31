package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmavrin/mock-http-server/internal/config"
)

func TestParseJSON(t *testing.T) {
	cases := []struct {
		it string

		config []byte

		expectedResult *config.Config
		expectedError  string
	}{
		{
			it: "it parse config",

			config: []byte(`
				{
					"prepareErrorPath": "/prepare-error",
					"handlers": [
						{
							"path": "/example-path",
							"method": "POST",
							"response": { "info": { "success": true }, "data": "post validate request success data" },
							"request": {
								"validate": {
									"example": "required,lowercase,min=5"
								},
								"match": { "example": "mock-request-body" }
							}
						}
					]
				}`),

			expectedResult: &config.Config{
				PrepareErrorPath: "/prepare-error",
				Handlers: []config.Handler{
					{
						Path:     "/example-path",
						Method:   config.Post,
						Response: []byte(`{ "info": { "success": true }, "data": "post validate request success data" }`),
						Request: &config.Request{
							Validate: map[string]any{
								"example": "required,lowercase,min=5",
							},
							Match: []byte(`{ "example": "mock-request-body" }`),
						},
					},
				},
			},
		},
		{
			it: "it parse other config",

			config: []byte(`
				{
					"prepareErrorPath": "/prepare-error",
					"handlers": [
						{
							"path": "/example-path",
							"method": "POST",
							"response": { "info": { "success": true }, "data": "post validate request success data" },
							"request": {
								"validate": {
									"example": {
										"mock": "value"
									}
								},
								"match": { "example": "mock-request-body" }
							}
						}
					]
				}`),

			expectedResult: &config.Config{
				PrepareErrorPath: "/prepare-error",
				Handlers: []config.Handler{
					{
						Path:     "/example-path",
						Method:   config.Post,
						Response: []byte(`{ "info": { "success": true }, "data": "post validate request success data" }`),
						Request: &config.Request{
							Validate: map[string]any{
								"example": map[string]any{
									"mock": "value",
								},
							},
							Match: []byte(`{ "example": "mock-request-body" }`),
						},
					},
				},
			},
		},
		{
			it: "it parse config with validate array",

			config: []byte(`
				{
					"prepareErrorPath": "/prepare-error",
					"handlers": [
						{
							"path": "/example-path",
							"method": "POST",
							"response": { "info": { "success": true }, "data": "post validate request success data" },
							"request": {
								"validate": [
									{
										"example": {
											"mock": "value"
										}
									}
								],
								"match": { "example": "mock-request-body" }
							}
						}
					]
				}`),

			expectedResult: &config.Config{
				PrepareErrorPath: "/prepare-error",
				Handlers: []config.Handler{
					{
						Path:     "/example-path",
						Method:   config.Post,
						Response: []byte(`{ "info": { "success": true }, "data": "post validate request success data" }`),
						Request: &config.Request{
							Validate: []any{
								map[string]any{
									"example": map[string]any{
										"mock": "value",
									},
								},
							},
							Match: []byte(`{ "example": "mock-request-body" }`),
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			data, err := config.Parse(config.JSON, tc.config)
			if err != nil {
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}

			assert.Equal(t, tc.expectedResult, data)
		})
	}
}
