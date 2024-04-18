package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	var csvDelimiter string
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

			output, err := run(file, format, csvDelimiter, pareto)
			if err != nil {
				return err
			}

			fmt.Print(output)

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "i", "", "Input JSON file")
	root.Flags().StringVarP(&format, "format", "f", "table", "Output format. One of [json table csv list fitness]")
	root.Flags().StringVarP(&csvDelimiter, "delim", "d", ",", "Column delimiter for CSV output")
	root.Flags().BoolVarP(&pareto, "pareto", "p", false, "Use pareto optimization criterion")

	return root
}

func run(file, format string, csvDelimiter string, pareto bool) (string, error) {
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

	b := strings.Builder{}
	switch format {
	case "json":
		enc := json.NewEncoder(&b)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "    ")
		err = enc.Encode(&solution)
		if err != nil {
			return "", err
		}
		b.WriteString(fmt.Sprintln(string(jsData)))

	case "table":
		for _, sol := range solution {
			b.WriteString(fmt.Sprintln(sol.ToTable()))
			b.WriteString(fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples))
			b.WriteString(fmt.Sprintln("------------------------------------------------------------"))
		}

	case "csv":
		for i, sol := range solution {
			b.WriteString(fmt.Sprint(sol.ToCSV(i, csvDelimiter)))
		}

	case "list":
		for _, sol := range solution {
			b.WriteString(fmt.Sprintln(sol.ToList()))
			b.WriteString(fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples))
			b.WriteString(fmt.Sprintln("------------------------------------------------------------"))
		}

	case "fitness":
		for _, sol := range solution {
			b.WriteString(fmt.Sprintf("(%d trips, %d samples)\n", sol.Fitness.Trips, sol.Fitness.Samples))
		}

	default:
		return "", fmt.Errorf("unknown format '%s'", format)
	}

	return b.String(), nil
}
