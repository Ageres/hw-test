package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageError(t *testing.T) {
	t.Run("simple error", func(t *testing.T) {
		err := NewSimpleSError("test error")
		require.NotNil(t, err)
		require.Equal(t, "test error", err.Error())
	})

	t.Run("error with template", func(t *testing.T) {
		err := NewSErrorWithTemplate("template %s %d", "test", 42)
		require.Equal(t, "template test 42", err.Error())
	})

	t.Run("error with empty template", func(t *testing.T) {
		err := NewSErrorWithTemplate("")
		require.Nil(t, err)
	})

	t.Run("error with message array", func(t *testing.T) {
		t.Run("non-empty messages", func(t *testing.T) {
			err := NewSErrorWithMsgArr([]string{"err1", "err2"})
			require.Equal(t, "err1; err2", err.Error())
		})

		t.Run("with empty messages", func(t *testing.T) {
			err := NewSErrorWithMsgArr([]string{"", "err1", ""})
			require.Equal(t, "err1", err.Error())
		})

		t.Run("all empty messages", func(t *testing.T) {
			err := NewSErrorWithMsgArr([]string{"", ""})
			require.Nil(t, err)
		})

		t.Run("empty array", func(t *testing.T) {
			err := NewSErrorWithMsgArr([]string{})
			require.Nil(t, err)
		})
	})

	t.Run("error with cause", func(t *testing.T) {
		cause := errors.New("original error")
		err := NewSErrorWithCause("wrapper: %v", cause)
		require.Equal(t, "wrapper: original error: original error", err.Error())

		var storageErr *StorageError
		ok := errors.As(err, &storageErr)
		require.True(t, ok)
		require.Equal(t, cause, storageErr.Cause)
		require.Equal(t, cause, errors.Unwrap(err))
	})

	t.Run("joinString function", func(t *testing.T) {
		require.Equal(t, "", joinString(nil))
		require.Equal(t, "", joinString([]string{}))
		require.Equal(t, "", joinString([]string{"", ""}))
		require.Equal(t, "a", joinString([]string{"a", ""}))
		require.Equal(t, "a; b", joinString([]string{"a", "", "b"}))
	})
}
