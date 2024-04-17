package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mlange-42/isso"
	"github.com/mlange-42/isso/fitness"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:\n$ isso problem.json")
		os.Exit(0)
	}

	file := os.Args[1]

	jsData, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	problem := isso.ProblemDef{}
	err = json.Unmarshal(jsData, &problem)
	if err != nil {
		log.Fatal(err)
	}

	p := isso.NewProblem(problem)

	s := isso.NewSolver(
		&fitness.TripsAndSamplesEvaluator{},
		&fitness.TripsThenSamples{},
	)

	if solution, ok := s.Solve(&p); ok {
		fmt.Fprintf(os.Stderr, "Found %d solution(s)\n\n", len(solution))

		jsData, err := json.MarshalIndent(&solution, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsData))
		return
	}
	fmt.Println("No solution found")
}
