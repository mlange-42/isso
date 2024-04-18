package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	out, err := run("../../data/problem.json", "fitness", ",", true)
	assert.Nil(t, err)
	assert.Equal(t, "(5 trips, 1826 samples)\n", out)

	_, err = run("../../data/problem.json", "json", ",", false)
	assert.Nil(t, err)

	_, err = run("../../data/problem.json", "table", ",", false)
	assert.Nil(t, err)

	_, err = run("../../data/problem.json", "csv", ",", false)
	assert.Nil(t, err)

	_, err = run("../../data/problem.json", "list", ",", false)
	assert.Nil(t, err)
}

func TestRootCommand(t *testing.T) {
	_ = RootCommand()
}
