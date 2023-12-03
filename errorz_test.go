package errorz_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/JFAexe/errorz"
)

func Test_Errorz(t *testing.T) {
	t.Parallel()

	var (
		err1 = errors.New("error1")
		err2 = errors.New("error2")
		err3 = errors.New("error3")
		err4 = fmt.Errorf("%w", err2)
		err5 = errors.Join(err1, err3)
	)

	t.Run("IsSingle", func(t *testing.T) {
		t.Parallel()

		require.False(t, errorz.IsSingle(err1))
		require.False(t, errorz.IsSingle(err5))
		require.True(t, errorz.IsSingle(err4))
	})

	t.Run("IsJoined", func(t *testing.T) {
		t.Parallel()

		require.False(t, errorz.IsJoined(err1))
		require.False(t, errorz.IsJoined(err4))
		require.True(t, errorz.IsJoined(err5))
	})

	t.Run("IsUnwrappable", func(t *testing.T) {
		t.Parallel()

		require.False(t, errorz.IsUnwrappable(err1))
		require.True(t, errorz.IsUnwrappable(err4))
		require.True(t, errorz.IsUnwrappable(err5))
	})

	t.Run("UnwrapAll", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, errorz.UnwrapAll(nil))
		require.Nil(t, errorz.UnwrapAll(err1))
		require.Equal(t, []error{err2}, errorz.UnwrapAll(err4))
		require.Equal(t, []error{err1, err3}, errorz.UnwrapAll(err5))
	})

	t.Run("IsMatching", func(t *testing.T) {
		t.Parallel()

		require.False(t, errorz.IsMatching(nil))
		require.False(t, errorz.IsMatching(err1))
		require.False(t, errorz.IsMatching(err1, nil))
		require.False(t, errorz.IsMatching(err5, err2))
		require.True(t, errorz.IsMatching(err1, err1))
		require.True(t, errorz.IsMatching(err4, err2))
	})

	t.Run("Matching", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, errorz.Matching(nil))
		require.Nil(t, errorz.Matching([]error{err1}))
		require.Equal(t, []error{err2}, errorz.Matching([]error{err2}, err1, err2))
		require.Equal(t, []error{err4}, errorz.Matching([]error{err4}, err1, err2))
		require.Equal(t, []error{err5}, errorz.Matching([]error{err5}, err1, err2))
		require.Equal(t, []error{err1, err3, err5}, errorz.Matching([]error{err1, err2, err3, err4, err5}, err1, err3))
	})

	t.Run("Allow", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, errorz.Allow(nil))
		require.Nil(t, errorz.Allow(err1))
		require.Nil(t, errorz.Allow(err4, err1))
		require.Equal(t, errorz.Allow(err1, err1), err1)
		require.Equal(t, errorz.Allow(err4, err2), err4)
		require.Equal(t, errorz.Allow(err5, err3), err5)
	})

	t.Run("Ignore", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, errorz.Ignore(nil))
		require.Nil(t, errorz.Ignore(err1))
		require.Nil(t, errorz.Ignore(err1, err1))
		require.Nil(t, errorz.Ignore(err5, err3))
		require.Equal(t, errorz.Ignore(err2, err1), err2)
		require.Equal(t, errorz.Ignore(err4, err3), err4)
	})
}
