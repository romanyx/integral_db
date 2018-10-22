package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	s := New()
	t.Log("Given initialized storage.")
	{
		t.Log("\t Test: 0\t When getting key is not defined, should return not found error.")
		{
			value, err := s.Get("key0")
			assert.Equal(t, ErrNotFound, err)
			assert.Nil(t, value)
		}

		t.Log("\t Test: 1\t When key is present should return value.")
		{
			k, v := "key1", "value"
			ctx := context.Background()
			s.Set(ctx, k, v)

			value, err := s.Get(k)

			assert.Equal(t, v, value)
			assert.Nil(t, err)
		}

		t.Log("\t Test: 2\t When context is done or cancel should delete value from storage.")
		{
			d := make(chan struct{})
			ctxDoneCall = func() {
				d <- struct{}{}
			}
			k, v := "key2", "value"
			ctx, cancel := context.WithCancel(context.Background())
			s.Set(ctx, k, v)
			cancel()
			<-d

			value, err := s.Get(k)

			assert.Equal(t, ErrNotFound, err)
			assert.Nil(t, value)
		}
		t.Log("\t Test: 3\t When same key is setting, should remove prevoius context done wait.")
		{
			k, v := "key3", "value"
			ctx, cancel := context.WithCancel(context.Background())
			s.Set(ctx, k, v)
			s.Set(context.Background(), k, v)
			cancel()

			value, err := s.Get(k)

			assert.Equal(t, v, value)
			assert.Nil(t, err)
		}
	}
}

func Benchmark_muxMap_Set(b *testing.B) {
	b.StopTimer()
	s := New()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		key := fmt.Sprintf("key_%d", n)
		b.StartTimer()

		s.Set(context.Background(), key, "value")
	}
}

func Benchmark_muxMap_Get(b *testing.B) {
	b.StopTimer()
	s := New()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		key := fmt.Sprintf("key_%d", n)
		s.Set(context.Background(), key, "value")
		b.StartTimer()

		s.Get(key)
	}
}

func Benchmark_muxMap_SetGet(b *testing.B) {
	b.StopTimer()
	s := New()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		key := fmt.Sprintf("key_%d", n)
		b.StartTimer()

		s.Set(context.Background(), key, "value")
		s.Get(key)
	}
}
