package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

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

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(1)
		c.Set("aaa", 100)

		c.Clear()
		res, flag := c.Get("aaa")
		require.Nil(t, res)
		require.False(t, flag)
	})
}
func TestCache_Add(t *testing.T) {
	t.Run("purge cash and set new cash", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Clear()
		res, flag := c.Get("aaa")
		require.Nil(t, res)
		require.False(t, flag)
		c.Set("bbb", 200)
		res, flag = c.Get("bbb")
		require.Equal(t, 200, res)
		require.True(t, flag)
		c.Clear()
	})

	t.Run("single value logic", func(t *testing.T) {
		c := NewCache(0)
		ok := c.Set("aaa", 100)
		require.False(t, ok)

		c = NewCache(1)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		val, flag := c.Get("aaa")
		require.False(t, flag)
		require.Nil(t, val)

		val, flag = c.Get("bbb")
		require.True(t, flag)
		require.Equal(t, 200, val)
	})

	t.Run("push logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		val, flag := c.Get("aaa")
		require.False(t, flag)
		require.Nil(t, val)
	})

	t.Run("logic of pushing out long-unused elements", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Get("aaa")
		c.Get("bbb")
		c.Get("aaa")
		c.Set("ddd", 400)
		val, flag := c.Get("ccc")
		require.False(t, flag)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Parallel()
	c := NewCache(10)
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
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000)))) //nolint
		}
	}()

	wg.Wait()
}
