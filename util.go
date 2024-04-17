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

func SolutionTable[F any](solution Solution[F]) string {
	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf("%10s %18s %6s %10s %10s %10s\n", "Subject", "Matrix", "Time", "Samples", "Reuse", "Target"),
	)

	for i, a := range solution.Actions {
		b.WriteString(
			fmt.Sprintf("%10s %18s %6d %10d %10s %10d", a.Subject, a.Matrix, a.Time, a.Samples, a.Reuse, a.TargetSamples),
		)
		if i < len(solution.Actions)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
