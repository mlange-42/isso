package isso_test

import (
	"testing"

	"github.com/mlange-42/isso"
)

func TestProblem(t *testing.T) {
	subjects := []string{
		"Pest 1",
		"Pest 2",
	}

	matrices := []isso.MatrixDef[string]{
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

	requirements := []isso.RequirementDef[string, string]{
		{
			Subject: "Pest 1",
			Matrix:  "shoots",
			Samples: 100,
			Times:   []int{3, 4, 5},
		},
	}

	p := isso.NewProblem(
		subjects,
		matrices,
		capacity,
		requirements,
	)

	_ = p
}
