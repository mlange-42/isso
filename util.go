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

// ToTable formats the solution as a table for printing.
func (s *Solution[F]) ToTable() string {
	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf("%10s %18s %6s %10s %10s %10s\n", "Subject", "Matrix", "Time", "Samples", "Reuse", "Target"),
	)

	for i, a := range s.Actions {
		b.WriteString(
			fmt.Sprintf("%10s %18s %6d %10d %10s %10d", a.Subject, a.Matrix, a.Time, a.Samples, a.Reuse, a.TargetSamples),
		)
		if i < len(s.Actions)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// ToCSV formats the solution as a CSV table.
func (s *Solution[F]) ToCSV(index int, sep string) string {
	b := strings.Builder{}

	if index <= 0 {
		if index >= 0 {
			b.WriteString(fmt.Sprintf("%s%s", "Solution", sep))
		}
		b.WriteString(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s\n", "Subject", sep, "Matrix", sep, "Time", sep, "Samples", sep, "Reuse", sep, "Target"))
	}

	for _, a := range s.Actions {
		if index >= 0 {
			b.WriteString(fmt.Sprintf("%d%s", index, sep))
		}
		b.WriteString(fmt.Sprintf("%s%s%s%s%d%s%d%s%s%s%d\n", a.Subject, sep, a.Matrix, sep, a.Time, sep, a.Samples, sep, a.Reuse, sep, a.TargetSamples))
	}
	return b.String()
}

type timeEntry struct {
	Subjects map[string]int
	Samples  int
}

// ToList formats the solution as list for printing.
func (s *Solution[F]) ToList() string {
	b := strings.Builder{}

	times := []map[string]timeEntry{}

	matrices := map[string]string{}
	for _, a := range s.Actions {
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

	for _, a := range s.Actions {
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
