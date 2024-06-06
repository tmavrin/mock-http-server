package validate_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Validate: map[string]string{
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
				Validate: map[string]string{
					"example": "required,lowercase,min=5",
				},
			},

			expectedError: "Key: '' Error:Field validation for '' failed on the 'min' tag", // TODO: check why return empty field name
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

			expectedError: `request does not match expected: { "example": "mock-request-body" } , actual: { "example": "mock" }`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			err := validate.RequestBody(tc.httpReq, tc.config)
			if err != nil {
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}
		})
	}
}
