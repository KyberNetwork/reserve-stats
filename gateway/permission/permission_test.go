package permission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractKeyID(t *testing.T) {
	var tests = []struct {
		header []string
		keyID  string
	}{
		{
			header: []string{
				`Signature keyId="Test",algorithm="rsa-sha256",signature="SjWJWbWN7i0wzBvtPl8rbASWz5xQW6mcJmn+ibttBqtifLN7Sazz6m79cNfwwb8DMJ5cou1s7uEGKKCs+FLEEaDV5lp7q25WqS+lavg7T8hc0GppauB6hbgEKTwblDHYGEtbGmtdHgVCk9SuS13F0hZ8FD0k/5OxEPXe5WozsbM="`,
			},
			keyID: "Test",
		},
		{
			header: []string{
				`keyId="Test",algorithm="rsa-sha256",signature="SjWJWbWN7i0wzBvtPl8rbASWz5xQW6mcJmn+ibttBqtifLN7Sazz6m79cNfwwb8DMJ5cou1s7uEGKKCs+FLEEaDV5lp7q25WqS+lavg7T8hc0GppauB6hbgEKTwblDHYGEtbGmtdHgVCk9SuS13F0hZ8FD0k/5OxEPXe5WozsbM="`,
			},
			keyID: "Test",
		},
		{
			header: []string{
				`keyId="Test",algorithm="rsa-sha256",signature="SjWJWbWN7i0wzBvtPl8rbASWz5xQW6mcJmn+ibttBqtifLN7Sazz6m79cNfwwb8DMJ5cou1s7uEGKKCs+FLEEaDV5lp7q25WqS+lavg7T8hc0GppauB6hbgEKTwblDHYGEtbGmtdHgVCk9SuS13F0hZ8FD0k/5OxEPXe5WozsbM="`,
			},
			keyID: "Test",
		},
		{
			header: []string{`a="keyId=",keyID="Test"`},
			keyID:  "Test",
		},
	}

	for _, tc := range tests {
		actualKeyID, err := extractKeyID(tc.header)
		require.NoError(t, err)
		assert.Equal(t, string(actualKeyID), tc.keyID)
	}
}
