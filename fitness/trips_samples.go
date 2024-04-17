package fitness

import (
	"github.com/mlange-42/isso"
)

type TripsAndSamplesFitness struct {
	Trips   int
	Samples int
}

type TripsAndSamplesEvaluator struct {
	times []int
}

func (e *TripsAndSamplesEvaluator) Evaluate(sol *isso.Actions) TripsAndSamplesFitness {
	for i := range e.times {
		e.times[i] = 0
	}
	samples := 0
	for _, a := range sol.Actions {
		for len(e.times) <= a.Time {
			e.times = append(e.times, 0)
		}
		e.times[a.Time] = 1
		if !a.IsReuse {
			samples += a.Samples
		}
	}
	trips := 0
	for _, t := range e.times {
		trips += t
	}

	return TripsAndSamplesFitness{
		trips,
		samples,
	}
}

type TripsAndSamplesComparator struct{}

func (e *TripsAndSamplesComparator) Less(a, b TripsAndSamplesFitness) bool {
	if b.Trips == 0 && b.Samples == 0 {
		return true
	}
	return a.Trips < b.Trips || (a.Trips == b.Trips && a.Samples < b.Samples)
}