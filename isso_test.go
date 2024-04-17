package isso_test

import (
	"testing"

	"github.com/mlange-42/isso"
	"github.com/mlange-42/isso/fitness"
)

func TestProblem(t *testing.T) {
	subjects := []string{
		"Pest 1",
		"Pest 2",
		"Pest 3",
		"Pest 4",
		"Pest 5",
		"Pest 6",
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
			Samples: 330,
			Times:   []int{2, 3, 4, 5},
		},
		{
			Subject: "Pest 2",
			Matrix:  "shoots",
			Samples: 419,
			Times:   []int{3, 4, 5, 6, 7},
		},
		{
			Subject: "Pest 3",
			Matrix:  "fruits",
			Samples: 970,
			Times:   []int{3, 4, 5, 6, 7, 9, 10, 11},
		},
		{
			Subject: "Pest 4",
			Matrix:  "fruits & shoots",
			Samples: 330,
			Times:   []int{8, 9, 10, 11},
		},
		{
			Subject: "Pest 5",
			Matrix:  "fruits & shoots",
			Samples: 1496,
			Times:   []int{3, 4, 5},
		},
		{
			Subject: "Pest 6",
			Matrix:  "fruits & shoots",
			Samples: 450,
			Times:   []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}

	p := isso.NewProblem(
		subjects,
		matrices,
		capacity,
		requirements,
	)

	s := isso.NewSolver[string, string, fitness.TripsAndSamplesFitness](
		&fitness.TripsAndSamplesEvaluator{},
		&fitness.TripsAndSamplesComparator{},
	)

	s.Solve(&p)
}
