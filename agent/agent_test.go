package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListenAndClose(t *testing.T) {
	err := Listen(Options{})
	assert.Nil(t, err)
	Close()
	_, err = os.Stat(pidFile)

	assert.True(t, os.IsNotExist(err))
	assert.Empty(t, pidFile)
}
