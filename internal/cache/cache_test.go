package cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) { //nolint:go-lint
	t.Run("empty cache", func(t *testing.T) {
		c, err := NewCache(10, "cache")
		require.NoError(t, err)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)

		err = c.Clear()
		require.NoError(t, err)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5, "cache")
		require.NoError(t, err)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		err = c.Clear()
		require.NoError(t, err)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c, err := NewCache(3, "cache")
		require.NoError(t, err)

		_ = c.Set("a", 100)
		_ = c.Set("b", 200)
		_ = c.Set("c", 300) // capacity reached

		val, ok := c.Get("a") // now "a" is in front of queue, "b" - in back
		require.True(t, ok)
		require.Equal(t, 100, val)

		_ = c.Set("d", 400)

		val, ok = c.Get("b")
		require.False(t, ok) // "b" has been pulled
		require.Equal(t, nil, val)

		_ = c.Set("e", 500)

		val, ok = c.Get("c")
		require.False(t, ok) // "c" has been pulled
		require.Equal(t, nil, val)

		_ = c.Set("f", 600)

		val, ok = c.Get("a")
		require.False(t, ok) // "a" has been pulled
		require.Equal(t, nil, val)

		// are "d", "e" and "f" in the place?
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 400, val)

		val, ok = c.Get("e")
		require.True(t, ok)
		require.Equal(t, 500, val)

		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 600, val)

		err = c.Clear()
		require.NoError(t, err)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c, err := NewCache(10, "cache")
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000)))) //nolint:go-lint
		}
	}()

	wg.Wait()
	err = c.Clear()
	require.NoError(t, err)
}
