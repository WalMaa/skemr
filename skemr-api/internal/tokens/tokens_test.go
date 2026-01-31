package tokens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	token, prefix, secret, err := GenerateToken(9, 32)

	require.NoError(t, err)
	require.NotEqual(t, 0, len(prefix))
	require.NotEqual(t, 0, len(secret))
	require.NotEqual(t, 0, len(token))

	verifier, err := HashSecret(secret, DefaultParams)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(verifier))
}
