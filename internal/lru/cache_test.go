package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	ch := make(chan string, 100)
	defer close(ch)

	go func() {
		for {
			<-ch
		}
	}()

	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10, ch)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5, ch)

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
	})

	t.Run("remove old without use", func(t *testing.T) {
		c := NewCache(3, ch)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("dddd", 100)
		require.False(t, wasInCache)
		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("remove old with use", func(t *testing.T) {
		c := NewCache(3, ch)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 100)
		require.False(t, wasInCache)
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 100, val)
		wasInCache = c.Set("dddd", 100)
		require.False(t, wasInCache)
		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	ch := make(chan string, 100)
	defer close(ch)

	go func() {
		for {
			<-ch
		}
	}()

	c := NewCache(10, ch)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(10000))))
		}
	}()

	wg.Wait()
}
