package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	// Test successful JWT creation and validation
	t.Run("Valid JWT", func(t *testing.T) {
		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		validatedUserID, err := ValidateJWT(token, tokenSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, validatedUserID)
	})

	// Test JWT signed with wrong secret is rejected
	t.Run("Wrong Secret", func(t *testing.T) {
		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		require.NoError(t, err)

		_, err = ValidateJWT(token, "wrong-secret")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	// Test expired JWT is rejected
	t.Run("Expired JWT", func(t *testing.T) {
		expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour) // Already expired
		require.NoError(t, err)

		_, err = ValidateJWT(expiredToken, tokenSecret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	// Test malformed JWT
	t.Run("Malformed JWT", func(t *testing.T) {
		_, err := ValidateJWT("not.a.jwt", tokenSecret)
		assert.Error(t, err)
	})

	// Test empty token
	t.Run("Empty Token", func(t *testing.T) {
		_, err := ValidateJWT("", tokenSecret)
		assert.Error(t, err)
	})

	// Test JWT that expires very soon
	t.Run("Short Expiry", func(t *testing.T) {
		shortExpiryToken, err := MakeJWT(userID, tokenSecret, time.Second)
		require.NoError(t, err)

		// Should work immediately
		validatedUserID, err := ValidateJWT(shortExpiryToken, tokenSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, validatedUserID)

		// Wait for expiration
		time.Sleep(time.Second + time.Millisecond*100)

		// Should now be expired
		_, err = ValidateJWT(shortExpiryToken, tokenSecret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})
}
