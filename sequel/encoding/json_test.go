package encoding

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type A struct {
	B *B
}

type B struct {
	ID int64
}

func TestJSON(t *testing.T) {
	a := A{}
	if a.B == nil {
		a.B = new(B)
	}
	require.NoError(t, JSONScanner(&a.B).Scan(nil))
	require.Nil(t, a.B)
}

type jsonRecord struct {
	data []byte
	err  error
}

func (r jsonRecord) MarshalJSON() ([]byte, error)     { return r.data, r.err }
func (r *jsonRecord) UnmarshalJSON(data []byte) error { r.data = data; return r.err }

func TestJSONValue(t *testing.T) {
	for _, tc := range []struct {
		input any
		want  any
	}{
		{nil, nil},
		{json.RawMessage(`{"id":42}`), []byte(`{"id":42}`)},
		{jsonRecord{data: []byte(`{"id":42}`)}, []byte(`{"id":42}`)},
		{map[string]int{"id": 42}, []byte("{\"id\":42}\n")},
	} {
		got, err := JSONValue(tc.input).Value()
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
	_, err := JSONValue(make(chan int)).Value()
	require.Error(t, err)
	sentinel := errors.New("marshal failed")
	_, err = JSONValue(jsonRecord{err: sentinel}).Value()
	require.ErrorIs(t, err, sentinel)
}

func TestJSONScanner(t *testing.T) {
	for _, input := range []any{`{"id":42}`, []byte(`{"id":42}`), json.RawMessage(`{"id":42}`)} {
		var got map[string]int
		scanner := JSONScanner(&got)
		require.NoError(t, scanner.Scan(input))
		require.Equal(t, map[string]int{"id": 42}, got)
		require.Error(t, scanner.Scan(42))
		require.Error(t, scanner.Scan(`{`))
		require.NoError(t, scanner.Scan(nil))
		require.Nil(t, got)
	}
	sentinel := errors.New("unmarshal failed")
	record := jsonRecord{data: []byte("old"), err: sentinel}
	require.ErrorIs(t, JSONScanner(&record).Scan(nil), sentinel)
	require.Nil(t, record.data)
}
