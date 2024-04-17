package isso

import (
	"fmt"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SolutionTable(solution []Action) string {
	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf("%10s %18s %6s %10s %10s %10s\n", "Subject", "Matrix", "Time", "Samples", "Reuse", "Target"),
	)

	for _, a := range solution {
		b.WriteString(
			fmt.Sprintf("%10s %18s %6d %10d %10s %10d\n", a.Subject, a.Matrix, a.Time, a.Samples, a.Reuse, a.TargetSamples),
		)
	}
	return b.String()
}
