package isso_test

import (
	"testing"

	"github.com/mlange-42/isso"
)

func TestProblem(t *testing.T) {
	matrices := []isso.MatrixDef{
		{Name: "fruits & shoots", CanReuse: []string{}},
		{Name: "fruits | shoots", CanReuse: []string{
			"fruits",
			"shoots",
			"fruits & shoots",
		}},
		{Name: "fruits", CanReuse: []string{
			"fruits & shoots",
		}},
		{Name: "shoots", CanReuse: []string{
			"fruits & shoots",
		}},
	}

	capacity := []int{
		150, 250, 400, 700, 600, 200, 50, 0, 150, 200, 150, 50,
	}

	p := isso.NewProblem(
		matrices,
		capacity,
	)

	_ = p
}
