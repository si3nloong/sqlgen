package strpool

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseResetsBuilder(t *testing.T) {
	b := AcquireString()
	require.Empty(t, b.String())
	b.WriteString("previous query")
	saved := b.String()
	ReleaseString(b)
	require.Zero(t, b.Len())
	next := AcquireString()
	defer ReleaseString(next)
	require.Empty(t, next.String())
	next.WriteString("next query")
	require.Equal(t, "previous query", saved)
	require.Equal(t, "next query", next.String())
}

func TestConcurrentAcquireRelease(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				b := AcquireString()
				if b.Len() != 0 {
					t.Error("acquired builder contains stale data")
				}
				b.WriteString("query")
				if b.String() != "query" {
					t.Error("builder contents changed")
				}
				ReleaseString(b)
			}
		}()
	}
	wg.Wait()
}
