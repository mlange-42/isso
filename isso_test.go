package isso_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/isso"
	"github.com/mlange-42/isso/fitness"
)

func TestDefaultProblem(t *testing.T) {
	subjects := []string{
		"Pest 1",
		"Pest 2",
		"Pest 3",
		"Pest 4",
		"Pest 5",
		"Pest 6",
	}

	matrices := []isso.Matrix{
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

	requirements := []isso.Requirement{
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
		isso.ProblemDef{
			Subjects:     subjects,
			Matrices:     matrices,
			Capacity:     capacity,
			Requirements: requirements,
		},
	)

	s := isso.NewSolver(
		&fitness.TripsAndSamplesEvaluator{},
		&fitness.TripsThenSamples{},
	)

	if solution, ok := s.Solve(&p); ok {
		fmt.Printf("Found %d solution(s)\n\n", len(solution))
		for _, sol := range solution {
			fmt.Println(isso.SolutionTable(sol))
			fmt.Println()
			fmt.Println(isso.SolutionList(sol))
			fmt.Println()
			fmt.Printf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples)
			fmt.Println()
			fmt.Println("------------------------------------------------------------")
			fmt.Println()
		}
		return
	}
	fmt.Println("No solution found")
}

func TestParetoProblem(t *testing.T) {
	subjects := []string{
		"Pest 1",
		"Pest 2",
	}

	matrices := []isso.Matrix{
		{Name: "fruits", CanReuse: []string{}},
	}

	capacity := []int{
		1000, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 1000,
	}

	requirements := []isso.Requirement{
		{
			Subject: "Pest 1",
			Matrix:  "fruits",
			Samples: 1000,
			Times:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			Subject: "Pest 2",
			Matrix:  "fruits",
			Samples: 1000,
			Times:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
	}

	p := isso.NewProblem(
		isso.ProblemDef{
			Subjects:     subjects,
			Matrices:     matrices,
			Capacity:     capacity,
			Requirements: requirements,
		},
	)

	s := isso.NewSolver(
		&fitness.TripsAndSamplesEvaluator{},
		&fitness.TripsSamplesPareto{},
	)

	if solution, ok := s.Solve(&p); ok {
		fmt.Printf("Found %d solution(s)\n\n", len(solution))
		for _, sol := range solution {
			fmt.Printf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples)
		}
		return
	}
	fmt.Println("No solution found")
}
