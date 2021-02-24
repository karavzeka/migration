package migration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerbose(t *testing.T) {
	log := Log{}
	assert.False(t, log.Verbose())
	log = Log{verbose: true}
	assert.True(t, log.Verbose())
}

func ExamplePrintf() {
	log := Log{}
	log.Printf("test %d %q", 12, "abc")
	// Output: test 12 "abc"
}
