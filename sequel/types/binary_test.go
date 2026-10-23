package types

import (
	"database/sql/driver"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

// base64Text should implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
type base64Text string

func (b base64Text) MarshalBinary() ([]byte, error) {
	// Implement your binary marshal logic here
	return []byte(base64.StdEncoding.EncodeToString([]byte(b))), nil
}

func (b *base64Text) UnmarshalBinary(data []byte) error {
	bytes, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	*b = base64Text(bytes)
	// Implement your binary unmarshal logic here
	return nil
}

func TestBinaryMarshaler(t *testing.T) {
	// Create an instance of base64Text
	v := base64Text(`abcd`) // Initialize your instance with appropriate values

	// Create a binaryMarshaler for base64Text
	bm := BinaryMarshaler(v)

	// Call the Value method of binaryMarshaler
	val, err := bm.Value()
	require.NoError(t, err)

	require.Equal(t, []byte("YWJjZA=="), val)

	// Assert that the returned value is of type driver.Value
	if _, ok := val.(driver.Value); !ok {
		t.Errorf("Value() did not return a value of type driver.Value")
	}
}

func TestBinaryUnmarshaler(t *testing.T) {
	// Create an instance of base64Text
	v := base64Text(``) // Initialize your instance with appropriate values

	// Create a binaryUnmarshaler for base64Text
	bum := BinaryUnmarshaler(&v)

	// Create a sample []byte data to be unmarshaled
	sampleData := []byte(`aGVsbG8gd29ybGQgIQ==`) // Provide a valid byte slice for testing

	// Call the Scan method of binaryUnmarshaler
	require.NoError(t, bum.Scan(sampleData))

	require.Equal(t, `hello world !`, string(v))
	// Assert that the base64Text instance has been properly unmarshaled
	// Compare the fields of yourInstance with the expected values
	// Add assertions based on your specific implementation
}

func TestBinaryUnmarshalerRejectsUnsupportedTypes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input any
	}{
		{"nil", nil},
		{"string", "YWJjZA=="},
		{"integer", int64(42)},
		{"boolean", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := base64Text("unchanged")
			err := BinaryUnmarshaler(&v).Scan(tc.input)
			require.EqualError(t, err, "sequel/types: binary must be []byte")
			require.Equal(t, base64Text("unchanged"), v)
		})
	}
}

func TestBinaryUnmarshalerPropagatesDecodeError(t *testing.T) {
	v := base64Text("unchanged")
	err := BinaryUnmarshaler(&v).Scan([]byte("!invalid base64!"))
	var decodeErr base64.CorruptInputError
	require.ErrorAs(t, err, &decodeErr)
	require.Equal(t, base64Text("unchanged"), v)
}
