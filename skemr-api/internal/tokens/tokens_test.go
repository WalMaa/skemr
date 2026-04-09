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

func TestVerifySecret(t *testing.T) {
	_, _, secret, err := GenerateToken(9, 32)
	require.NoError(t, err)

	verifier, err := HashSecret(secret, DefaultParams)
	require.NoError(t, err)

	ok, err := VerifySecret(secret, verifier)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = VerifySecret("wrongsecret", verifier)
	require.NoError(t, err)
	require.False(t, ok)
}
