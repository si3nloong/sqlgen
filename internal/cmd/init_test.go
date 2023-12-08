package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/require"
)

func TestNoInterruptError(t *testing.T) {
	var err = errors.New("done")
	var alias = io.EOF
	require.Equal(t, nil, noInterruptError(terminal.InterruptErr))
	require.Equal(t, err, noInterruptError(err))
	require.Equal(t, io.EOF, noInterruptError(alias))
}
