package fitness

import (
	"cmp"

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

type TripsThenSamples struct{}

func (e *TripsThenSamples) Compare(a, b TripsAndSamplesFitness) int {
	if b.Trips == 0 && b.Samples == 0 {
		return -1
	}
	if a.Trips < b.Trips {
		return -1
	}
	if a.Trips == b.Trips {
		return cmp.Compare(a.Samples, b.Samples)
	}
	return 1
}

func (e *TripsThenSamples) IsPareto() bool {
	return false
}

type TripsSamplesPareto struct{}

func (e *TripsSamplesPareto) Compare(a, b TripsAndSamplesFitness) int {
	if b.Trips == 0 && b.Samples == 0 {
		return -1
	}
	if a.Trips == b.Trips {
		return cmp.Compare(a.Samples, b.Samples)
	}
	if a.Samples == b.Samples {
		return cmp.Compare(a.Trips, b.Trips)
	}
	if a.Trips < b.Trips && a.Samples < b.Samples {
		return -1
	}
	if a.Trips > b.Trips && a.Samples > b.Samples {
		return 1
	}
	return 0
}

func (e *TripsSamplesPareto) IsPareto() bool {
	return true
}
