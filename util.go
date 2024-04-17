package isso

import (
	"fmt"
	"sort"
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

type timeEntry struct {
	Samples  int
	Subjects map[string]int
}

func SolutionList[F any](solution Solution[F]) string {
	b := strings.Builder{}

	times := []map[string]timeEntry{}

	matrices := map[string]string{}
	for _, a := range solution.Actions {
		for len(times) <= a.Time {
			times = append(times, map[string]timeEntry{})
		}

		if a.Reuse == "" {
			t := times[a.Time]
			t[a.Matrix] = timeEntry{
				Samples:  a.Samples,
				Subjects: map[string]int{a.Subject: a.Samples},
			}
			matrices[a.Subject] = a.Matrix
		}
	}

	for _, a := range solution.Actions {
		if a.Reuse != "" {
			t := times[a.Time]
			matrix := matrices[a.Reuse]
			t[matrix].Subjects[a.Subject] = a.Samples
		}
	}

	lines := []string{}
	for i, t := range times {
		if len(t) == 0 {
			continue
		}
		first := true
		for matrix, entry := range t {
			if first {
				lines = append(lines, fmt.Sprintf("Time = %2d: %4d x %-16s", i, entry.Samples, matrix))
			} else {
				lines = append(lines, fmt.Sprintf("           %4d x %-16s", entry.Samples, matrix))
			}
			keys := make([]string, 0, len(entry.Subjects))
			for k := range entry.Subjects {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, sub := range keys {
				lines = append(lines, fmt.Sprintf("               %4d x %-16s", entry.Subjects[sub], sub))
			}

			first = false
		}
	}
	b.WriteString(strings.Join(lines, "\n"))

	return b.String()
}
