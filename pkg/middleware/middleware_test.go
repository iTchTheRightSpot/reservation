package middleware

import (
	"os"
	"testing"
)

func TestMiddleware(t *testing.T) {
	if err := os.Chdir("../../"); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Parallel()

	t.Run("should reject request invalid cookie", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should reject request invalid jwt", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should accept request valid jwt", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should reject request. access denied valid jwt but not valid role", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should accept request valid jwt matching role", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should accept request & refresh token", func(t *testing.T) {
		t.Parallel()
	})
}
