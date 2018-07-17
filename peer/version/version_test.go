

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()
	assert.NoError(t, cmd.Execute(), "expected version command to succeed")
}
