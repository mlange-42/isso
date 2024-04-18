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
	var pareto bool

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
			if file == "" {
				_ = cmd.Help()
				return nil
			}

			output, err := run(file, format, pareto)
			if err != nil {
				return err
			}

			fmt.Println(output)

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "i", "", "Input JSON file")
	root.Flags().StringVarP(&format, "format", "f", "table", "Output format. One of [json table list fitness]")
	root.Flags().BoolVarP(&pareto, "pareto", "p", false, "Use pareto optimization criterion")

	return root
}

func run(file, format string, pareto bool) (string, error) {
	jsData, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	problem := isso.ProblemDef{}
	err = json.Unmarshal(jsData, &problem)
	if err != nil {
		return "", err
	}

	p := isso.NewProblem(problem)

	var comp isso.Comparator[fitness.TripsAndSamplesFitness]
	if pareto {
		comp = &fitness.TripsSamplesPareto{}
	} else {
		comp = &fitness.TripsThenSamples{}
	}

	s := isso.NewSolver(
		&fitness.TripsAndSamplesEvaluator{},
		comp,
	)
	solution, ok := s.Solve(&p)
	if !ok {
		fmt.Println("No solution found")
		return "", nil
	}

	fmt.Fprintf(os.Stderr, "Found %d solution(s)\n\n", len(solution))

	output := ""
	switch format {
	case "json":
		jsData, err = json.MarshalIndent(&solution, "", "    ")
		if err != nil {
			return "", err
		}
		fmt.Println(string(jsData))

	case "table":
		for _, sol := range solution {
			output += fmt.Sprintln(isso.SolutionTable(sol))
			output += fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples)
			output += fmt.Sprintln("------------------------------------------------------------")
		}
	case "list":
		for _, sol := range solution {
			output += fmt.Sprintln(isso.SolutionList(sol))
			output += fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples)
			output += fmt.Sprintln("------------------------------------------------------------")
		}
	case "fitness":
		for _, sol := range solution {
			output += fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples)
		}

	default:
		return "", fmt.Errorf("unknown format '%s'", format)
	}

	return output, nil
}
