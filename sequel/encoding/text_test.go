package encoding

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type textRecord struct {
	data string
	err  error
}

func (r *textRecord) UnmarshalText(b []byte) error {
	if r.err != nil {
		return r.err
	}
	r.data = string(b)
	return nil
}
func (r textRecord) MarshalText() ([]byte, error) { return []byte(r.data), r.err }

func TestTextScanner(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	for _, tc := range []struct {
		input any
		want  string
	}{
		{"hello", "hello"}, {[]byte("world"), "world"}, {true, "\x01"}, {false, "\x00"}, {int64(-42), "-42"}, {1.25, "1.25"}, {now, "2026-01-02T03:04:05Z"},
	} {
		for _, pointer := range []bool{false, true} {
			value := textRecord{}
			ptr := &value
			var scanner sql.Scanner
			if pointer {
				scanner = TextScanner[textRecord](&ptr)
			} else {
				scanner = TextScanner[textRecord](&value)
			}
			require.NoError(t, scanner.Scan(tc.input))
			require.Equal(t, tc.want, value.data)
			require.Error(t, scanner.Scan(struct{}{}))
			require.Error(t, scanner.Scan(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)))
			require.Equal(t, tc.want, value.data)
			require.NoError(t, scanner.Scan(nil))
			if pointer {
				require.Nil(t, ptr)
			} else {
				require.Equal(t, textRecord{}, value)
			}
		}
	}
}

func TestTextErrorsAndValue(t *testing.T) {
	sentinel := errors.New("text conversion failed")
	for _, pointer := range []bool{false, true} {
		value := textRecord{err: sentinel}
		ptr := &value
		var scanner sql.Scanner
		if pointer {
			scanner = TextScanner[textRecord](&ptr)
		} else {
			scanner = TextScanner[textRecord](&value)
		}
		require.ErrorIs(t, scanner.Scan("input"), sentinel)
	}
	value, err := TextValue(textRecord{data: "encoded"}).Value()
	require.NoError(t, err)
	require.Equal(t, []byte("encoded"), value)
	_, err = TextValue(textRecord{err: sentinel}).Value()
	require.ErrorIs(t, err, sentinel)
}
