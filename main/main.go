package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mlange-42/isso"
	"github.com/mlange-42/isso/fitness"
	"github.com/spf13/cobra"
)

func main() {
	if err := RootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// RootCommand sets up the CLI
func RootCommand() *cobra.Command {
	var format string
	var file string

	root := &cobra.Command{
		Use:           "isso",
		Short:         "isso -- Iterative Sampling Schedule Optimization",
		Long:          `isso -- Iterative Sampling Schedule Optimization`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			jsData, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			problem := isso.ProblemDef{}
			err = json.Unmarshal(jsData, &problem)
			if err != nil {
				return err
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
					return err
				}
				fmt.Println(string(jsData))
				return nil
			}
			fmt.Println("No solution found")

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "i", "", "Input JSON file")
	root.Flags().StringVarP(&format, "format", "f", "table", "Output format. One of [table json]")

	return root
}
