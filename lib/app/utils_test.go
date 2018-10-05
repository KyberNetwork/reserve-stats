package app

import "testing"

func TestJoinEnvVar(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		envVar string
		result string
	}{
		{
			name:   "prefix without traling _",
			prefix: "PREFIX",
			envVar: "ENV_VAR",
			result: "PREFIX_ENV_VAR",
		},
		{
			name:   "prefix with traling _",
			prefix: "PREFIX_",
			envVar: "ENV_VAR",
			result: "PREFIX_ENV_VAR",
		},
		{
			name:   "envVar with leading _",
			prefix: "PREFIX_",
			envVar: "_ENV_VAR",
			result: "PREFIX_ENV_VAR",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			joined := joinEnvVar(testCase.prefix, testCase.envVar)
			if joined != testCase.result {
				t.Errorf("got unexpected result for prefix: %s, envVar: %s, expected: %s, got: %s",
					testCase.prefix,
					testCase.envVar,
					testCase.result,
					joined,
				)
			}
		})
	}
}
