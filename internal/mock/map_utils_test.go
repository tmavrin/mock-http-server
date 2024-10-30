package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindValue(t *testing.T) {
	cases := []struct {
		it string

		source map[string]any
		key    string

		expectedResult any
	}{
		{
			it: "find value in simple 1 level json",

			source: map[string]any{
				"data_prop": "prop_value",
			},
			key: "data_prop",

			expectedResult: "prop_value",
		},
		{
			it: "find value in simple 6 level json",

			source: map[string]any{
				"data_prop_1": map[string]any{
					"data_prop_2": map[string]any{
						"data_prop_3": map[string]any{
							"data_prop_4": map[string]any{
								"data_prop_5": map[string]any{
									"data_prop_6": "prop_value",
								},
							},
						},
					},
				},
			},
			key:            "data_prop_1.data_prop_2.data_prop_3.data_prop_4.data_prop_5.data_prop_6",
			expectedResult: "prop_value",
		},
		{
			it: "return nil if not found",

			source: map[string]any{
				"data_prop_1": map[string]any{
					"data_prop_2": map[string]any{
						"data_prop_3": map[string]any{
							"data_prop_4": map[string]any{
								"data_prop_5": map[string]any{
									"data_prop_6": "prop_value",
								},
							},
						},
					},
				},
			},
			key:            "data_prop_1.data_prop_2.data_prop_10.data_prop_4.data_prop_5.data_prop_6",
			expectedResult: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {

			result := findValue(tc.source, tc.key)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestSetValue(t *testing.T) {
	cases := []struct {
		it string

		dest  map[string]any
		key   string
		value any

		expectedResult any
	}{
		{
			it: "find value in simple 1 level json",

			dest: map[string]any{
				"data_prop": "prop_value",
			},
			key:   "data_prop",
			value: "echo_value",

			expectedResult: map[string]any{
				"data_prop": "echo_value",
			},
		},
		{
			it: "find value in simple 6 level json",

			dest: map[string]any{
				"data_prop_1": map[string]any{
					"data_prop_2": map[string]any{
						"data_prop_3": map[string]any{
							"data_prop_4": map[string]any{
								"data_prop_5": map[string]any{
									"data_prop_6": "prop_value",
								},
							},
						},
					},
				},
			},
			key:   "data_prop_1.data_prop_2.data_prop_3.data_prop_4.echo_prop",
			value: "echo_value",

			expectedResult: map[string]any{
				"data_prop_1": map[string]any{
					"data_prop_2": map[string]any{
						"data_prop_3": map[string]any{
							"data_prop_4": map[string]any{
								"echo_prop": "echo_value",
								"data_prop_5": map[string]any{
									"data_prop_6": "prop_value",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {

			setValue(tc.dest, tc.key, tc.value)

			assert.Equal(t, tc.expectedResult, tc.dest)
		})
	}
}
