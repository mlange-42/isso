package isso_test

import (
	"testing"

	"github.com/mlange-42/isso"
)

func TestProblem(t *testing.T) {
	p := isso.NewProblem(
		[]string{
			"fruits & shoots",
			"fruits | shoots",
			"fruits",
			"shoots",
		},
		[]int{
			150, 250, 400, 700, 600, 200, 50, 0, 150, 200, 150, 50,
		},
	)

	_ = p
}
