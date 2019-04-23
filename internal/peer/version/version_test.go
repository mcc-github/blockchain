

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()
	assert.NoError(t, cmd.Execute(), "expected version command to succeed")
}

func TestCmdWithTrailingArgs(t *testing.T) {
	cmd := Cmd()
	args := []string{"trailingargs"}
	cmd.SetArgs(args)
	assert.EqualError(t, cmd.Execute(), "trailing args detected")
}
