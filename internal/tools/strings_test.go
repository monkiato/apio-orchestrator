package tools

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestIsAlpha(t *testing.T) {
	t.Run("should be alphabetic string only", func(t *testing.T) {
		assert.Equal(t, IsAlpha("testingAlphabeticStringOnly"), true)
		assert.Equal(t, IsAlpha("abcdefghijklmnopqrstuvwxyz"), true)
		assert.Equal(t, IsAlpha("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), true)
	})
	t.Run("should NOT be alphabetic", func(t *testing.T) {
		assert.Equal(t, IsAlpha("testing-Alphabetic-String"), false)
		assert.Equal(t, IsAlpha("value with spaces"), false)
		assert.Equal(t, IsAlpha("123456"), false)
	})
}

func TestIsAlphanumeric(t *testing.T) {
	t.Run("should be alphanumeric only", func(t *testing.T) {
		assert.Equal(t, IsAlphanumeric("testingAlphanumeric1234567890"), true)
		assert.Equal(t, IsAlphanumeric("abcdefghijklmnopqrstuvwxyz1234567890"), true)
		assert.Equal(t, IsAlphanumeric("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"), true)
	})
	t.Run("should NOT be alphanumeric", func(t *testing.T) {
		assert.Equal(t, IsAlphanumeric("testing-Alphanumeric-123456789"), false)
		assert.Equal(t, IsAlphanumeric("value with spaces 1234567890"), false)
		assert.Equal(t, IsAlphanumeric("&123456&"), false)
		assert.Equal(t, IsAlphanumeric("!123456"), false)
		assert.Equal(t, IsAlphanumeric("@123456"), false)
		assert.Equal(t, IsAlphanumeric("[123456]"), false)
	})
}

func TestIsValidFormat(t *testing.T) {
	t.Run("should use limited chars", func(t *testing.T) {
		assert.Equal(t, IsValidFormat("my.new_email@gmail.com", "@._abcdefghijklmnopqrstuvwxyz1234567890"), true)
		assert.Equal(t, IsValidFormat("invalid$$@gmail.com", "@._abcdefghijklmnopqrstuvwxyz1234567890"), false)
		assert.Equal(t, IsValidFormat("f045ce789d", "1234567890abcdef"), true)
		assert.Equal(t, IsValidFormat("invalid-hex", "1234567890abcdef"), false)
	})
}