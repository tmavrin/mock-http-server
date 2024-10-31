package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmavrin/mock-http-server/internal/config"
	"github.com/tmavrin/mock-http-server/internal/handler"
)

func TestHandler(t *testing.T) {

	cases := []struct {
		it string

		request map[string]any

		config config.Handler

		expectedResult string
		expectedStatus int
	}{
		{
			it: "validate and match",

			request: map[string]any{"example": "mock-request-body"},

			config: config.Handler{
				Method:   http.MethodPost,
				Path:     "/example",
				Response: []byte(`{"info":{"success":true},"data":"post validate request success data"}`),
				Request: &config.Request{
					Validate: map[string]any{
						"example": "required,lowercase,min=5",
					},
					Match: []byte(`{"example":"mock-request-body"}`),
				},
			},

			expectedStatus: http.StatusOK,
			expectedResult: `{"info":{"success":true},"data":"post validate request success data"}`,
		},
		{
			it: "validate, match and echo request",

			request: map[string]any{
				"example":  "mock-request-body",
				"testEcho": "echoValue",
			},

			config: config.Handler{
				Method:   http.MethodPost,
				Path:     "/example",
				Response: []byte(`{"info":{"success":true},"data":"post validate request success data"}`),
				Request: &config.Request{
					Validate: map[string]any{
						"example": "required,lowercase,min=5",
					},
					Match: []byte(`{"example":"mock-request-body"}`),
				},
				ResponseEcho: map[string]string{
					"testEcho": "data.testEcho",
				},
			},

			expectedStatus: http.StatusOK,
			expectedResult: `{"info":{"success":true},"data":{"testEcho":"echoValue"}}`,
		},

		{
			it: "validate query params",

			config: config.Handler{
				Method:   http.MethodGet,
				Path:     "/get?test=mock-value",
				Response: []byte(`{"info":{"success":true},"data":"get request success data"}`),
				QueryParams: map[string]config.QueryParam{
					"test": {
						Required: true,
						Value:    "mock-value",
					},
				},
			},

			expectedStatus: http.StatusOK,
			expectedResult: `{"info":{"success":true},"data":"get request success data"}`,
		},
		{
			it: "validate non required query params",

			config: config.Handler{
				Method:   http.MethodGet,
				Path:     "/get",
				Response: []byte(`{"info":{"success":true},"data":"get request success data"}`),
				QueryParams: map[string]config.QueryParam{
					"test": {
						Required: false,
						Value:    "mock-value",
					},
				},
			},

			expectedStatus: http.StatusOK,
			expectedResult: `{"info":{"success":true},"data":"get request success data"}`,
		},
		{
			it: "validate bad query params",

			config: config.Handler{
				Method:   http.MethodGet,
				Path:     "/get?test=invalid-value",
				Response: []byte(`{"info":{"success":true},"data":"get request success data"}`),
				QueryParams: map[string]config.QueryParam{
					"test": {
						Required: true,
						Value:    "mock-value",
					},
				},
			},

			expectedStatus: http.StatusBadRequest,
			expectedResult: `"query param test value 'invalid-value' does not match 'mock-value'"`,
		},
		{
			it: "pass null as value on echo request without sent request property",

			request: map[string]any{
				"example": "mock-request-body",
			},

			config: config.Handler{
				Method:   http.MethodPost,
				Path:     "/example",
				Response: []byte(`{"info":{"success":true},"data":"post validate request success data"}`),
				Request: &config.Request{
					Validate: map[string]any{
						"example": "required,lowercase,min=5",
					},
					Match: []byte(`{"example":"mock-request-body"}`),
				},
				ResponseEcho: map[string]string{
					"testEcho": "data.testEcho",
				},
			},

			expectedStatus: http.StatusOK,
			expectedResult: `{"info":{"success":true},"data":{"testEcho":null}}`,
		},

		{
			it: "fail on validate",

			request: map[string]any{"example": "mock"},

			config: config.Handler{
				Method:   http.MethodPost,
				Path:     "/example",
				Response: []byte(`{"info":{"success":true},"data":"post validate request success data"}`),
				Request: &config.Request{
					Validate: map[string]any{
						"example": "required,lowercase,min=5",
					},
				},
			},

			expectedStatus: http.StatusBadRequest,
			expectedResult: `"Key: 'example' Error:Field validation for 'example' failed on the 'min' tag"`,
		},

		{
			it: "fail on match",

			request: map[string]any{"example": "mock"},

			config: config.Handler{
				Method:   http.MethodPost,
				Path:     "/example",
				Response: []byte(`{"info":{"success":true},"data":"post validate request success data"}`),
				Request: &config.Request{
					Match: []byte(`{"example":"mock-request-body"}`),
				},
			},

			expectedStatus: http.StatusBadRequest,
			expectedResult: `"request does not match expected: {\"example\":\"mock-request-body\"} , actual: {\"example\":\"mock\"}"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {

			data, err := json.Marshal(tc.request)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.config.Method.String(), tc.config.Path, bytes.NewReader(data))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.Handler(tc.config).ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.JSONEq(t, tc.expectedResult, rr.Body.String())
		})
	}
}
