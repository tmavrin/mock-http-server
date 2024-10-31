package validate_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmavrin/mock-http-server/internal/config"
	"github.com/tmavrin/mock-http-server/internal/validate"
)

func TestRequestBody(t *testing.T) {
	cases := []struct {
		it string

		httpReq *http.Request
		config  *config.Request

		expectedError string
	}{
		{
			it: "it validate request body by fields",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{
								"template":"mock-template",
								"example":"mock-example"
							}`,
						),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": "required,lowercase,min=5",
				},
			},
		},
		{
			it: "it return error on invalid request body by fields",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{
								"template":"mock-template",
								"example":"mock"
							}`,
						),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": "required,lowercase,min=5",
				},
			},

			expectedError: "Key: 'example' Error:Field validation for 'example' failed on the 'min' tag",
		},
		{
			it: "it validate request body by exact match",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{ "example": "mock-request-body" }`),
					),
				),
			},

			config: &config.Request{
				Match: []byte(`{ "example": "mock-request-body" }`),
			},
		},

		{
			it: "it return error on request body by exact match",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{ "example": "mock" }`),
					),
				),
			},

			config: &config.Request{
				Match: []byte(`{ "example": "mock-request-body" }`),
			},

			expectedError: `request does not match expected: { "example": "mock-request-body" } , actual: {"example":"mock"}`,
		},

		{
			it: "it processes object",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{ "example": { "mock": "value" } }`),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": map[string]any{
						"mock": "required",
					},
				},
			},
		},

		{
			it: "it processes object validation",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{ "example": { "mock": null } }`),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": map[string]any{
						"mock": "required",
					},
				},
			},

			expectedError: "Key: 'mock' Error:Field validation for 'mock' failed on the 'required' tag",
		},

		{
			it: "it processes multi layer object validation",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`{ "example": { "mock": { "mock-2": { "mock-3": null } } } }`),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": map[string]any{
						"mock": map[string]any{
							"mock-2": map[string]any{
								"mock-3": "required",
							},
						},
					},
				},
			},

			expectedError: "Key: 'mock-3' Error:Field validation for 'mock-3' failed on the 'required' tag",
		},

		{
			it: "it processes array validation",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`[{ "example": { "mock": { "mock-2": { "mock-3": null } } } }]`),
					),
				),
			},

			config: &config.Request{
				Validate: []any{
					map[string]any{
						"example": map[string]any{
							"mock": map[string]any{
								"mock-2": map[string]any{
									"mock-3": "required",
								},
							},
						},
					},
				},
			},

			expectedError: "Key: 'mock-3' Error:Field validation for 'mock-3' failed on the 'required' tag",
		},
		{
			it: "it processes object validation with array as property",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer([]byte(`[{"example": {"mock": [{"mock-2": {"mock-3": null}}]} }]`)),
				),
			},

			config: &config.Request{
				Validate: []any{
					map[string]any{
						"example": map[string]any{
							"mock": []any{
								map[string]any{
									"mock-2": map[string]any{
										"mock-3": "required",
									},
								},
							},
						},
					},
				},
			},

			expectedError: "Key: 'mock-3' Error:Field validation for 'mock-3' failed on the 'required' tag",
		},
		{
			it: "it processes object validation with array as property",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer([]byte(`[{},{"example": {"mock": [{"mock-2": {"mock-3": null}}]} }]`)),
				),
			},

			config: &config.Request{
				Validate: []any{
					map[string]any{},
					map[string]any{
						"example": map[string]any{
							"mock": []any{
								map[string]any{
									"mock-2-1": map[string]any{
										"mock-2-3": "required",
									},
								},
							},
						},
					},
				},
			},

			expectedError: "Key: 'mock-2-3' Error:Field validation for 'mock-2-3' failed on the 'required' tag",
		},
		{
			it: "it returns null if no validation",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer([]byte(`[{"example": {"mock": [{"mock-2": {"mock-3": null}}]} }]`)),
				),
			},

			config: &config.Request{},
		},

		{
			it: "it returns type mismatch",

			httpReq: &http.Request{
				Body: io.NopCloser(
					bytes.NewBuffer(
						[]byte(`[{ "example": { "mock": { "mock-2": { "mock-3": null } } } }]`),
					),
				),
			},

			config: &config.Request{
				Validate: map[string]any{
					"example": map[string]any{
						"mock": map[string]any{
							"mock-2": map[string]any{
								"mock-3": "required",
							},
						},
					},
				},
			},

			expectedError: "config request validation is defined as an object, but sent request is not an object",
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			var httpBody any

			err := json.NewDecoder(tc.httpReq.Body).Decode(&httpBody)

			require.NoError(t, err)

			err = validate.RequestBody(httpBody, tc.config)
			if err != nil {
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}
		})
	}
}
