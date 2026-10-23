package encoding

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type recordingScanner struct {
	input any
	err   error
}

func (r *recordingScanner) Scan(input any) error { r.input = input; return r.err }

func TestPtrScannerDelegates(t *testing.T) {
	sentinel := errors.New("scan failed")
	for _, input := range []any{"text", []byte("bytes"), true, int64(42), 1.5, time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)} {
		value := &recordingScanner{}
		scanner := PtrScanner[recordingScanner](&value)
		require.NoError(t, scanner.Scan(input))
		require.Equal(t, input, value.input)
		value.err = sentinel
		require.ErrorIs(t, scanner.Scan(input), sentinel)
		require.Error(t, scanner.Scan(struct{}{}))
		require.NoError(t, scanner.Scan(nil))
		require.Nil(t, value)
	}
}
