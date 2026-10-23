package encoding

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type scanCase[T any] struct {
	input any
	want  T
}

// Check both destination forms, NULL resets, and preservation on failed scans.
func checkScanner[T any](t *testing.T, direct func(*T) sql.Scanner, nullable func(**T) sql.Scanner, cases []scanCase[T], invalid []any) {
	t.Helper()
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%T/%v", tc.input, tc.input), func(t *testing.T) {
			var value T
			var ptr *T
			for _, scanner := range []sql.Scanner{direct(&value), nullable(&ptr)} {
				require.NoError(t, scanner.Scan(tc.input))
			}
			require.Equal(t, tc.want, value)
			require.NotNil(t, ptr)
			require.Equal(t, tc.want, *ptr)
			for _, input := range invalid {
				require.Error(t, direct(&value).Scan(input))
				require.Error(t, nullable(&ptr).Scan(input))
				require.Equal(t, tc.want, value)
				require.Equal(t, tc.want, *ptr)
			}
			require.NoError(t, direct(&value).Scan(nil))
			require.NoError(t, nullable(&ptr).Scan(nil))
			var zero T
			require.Equal(t, zero, value)
			require.Nil(t, ptr)
		})
	}
}
