package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopyMap(t *testing.T) {
	m1 := map[string]interface{}{
		"EUR": "Something",
		"USD": map[string]interface{}{
			"balanceId": 123,
		},
	}

	m2 := CopyMap(m1)

	m1["EUR"] = "bumble"
	delete(m1, "USD")

	require.Equal(t, map[string]interface{}{"EUR": "bumble"}, m1)
	require.Equal(t, map[string]interface{}{
		"EUR": "Something",
		"USD": map[string]interface{}{
			"balanceId": 123,
		},
	}, m2)
}
