package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	out, err := run("../data/problem.json", "fitness", true)
	assert.Nil(t, err)
	assert.Equal(t, "(5 trips, 1826 samples)\n", out)
}
