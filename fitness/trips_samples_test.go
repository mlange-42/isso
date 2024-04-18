package fitness_test

import (
	"testing"

	"github.com/mlange-42/isso"
	"github.com/mlange-42/isso/fitness"
	"github.com/stretchr/testify/assert"
)

type f = fitness.TripsAndSamplesFitness

func TestTripsAndSamplesEvaluator(t *testing.T) {
	solution := isso.Actions{Actions: []isso.ActionDef{
		{
			Subject: 1,
			Samples: 100,
			Time:    2,
		},
		{
			Subject: 2,
			Samples: 50,
			Time:    4,
		},
		{
			Subject: 3,
			Samples: 50,
			Time:    4,
			Reuse:   2,
			IsReuse: true,
		},
	}}

	eval := fitness.TripsAndSamplesEvaluator{}
	fit := eval.Evaluate(&solution)

	assert.Equal(t, f{Trips: 2, Samples: 150}, fit)
}

func TestTripsThenSamples(t *testing.T) {
	comp := fitness.TripsThenSamples{}

	assert.False(t, comp.IsPareto())

	assert.Equal(t, -1, comp.Compare(
		f{Trips: 1, Samples: 1000},
		f{Trips: 2, Samples: 100},
	))

	assert.Equal(t, -1, comp.Compare(
		f{Trips: 1, Samples: 100},
		f{Trips: 2, Samples: 100},
	))

	assert.Equal(t, 0, comp.Compare(
		f{Trips: 1, Samples: 100},
		f{Trips: 1, Samples: 100},
	))

	assert.Equal(t, 1, comp.Compare(
		f{Trips: 2, Samples: 100},
		f{Trips: 1, Samples: 1000},
	))

	assert.Equal(t, 1, comp.Compare(
		f{Trips: 1, Samples: 101},
		f{Trips: 1, Samples: 100},
	))
}

func TestTripsSamplesPareto(t *testing.T) {
	comp := fitness.TripsSamplesPareto{}

	assert.True(t, comp.IsPareto())

	assert.Equal(t, -1, comp.Compare(
		f{Trips: 1, Samples: 100},
		f{Trips: 2, Samples: 100},
	))
	assert.Equal(t, -1, comp.Compare(
		f{Trips: 1, Samples: 100},
		f{Trips: 1, Samples: 200},
	))

	assert.Equal(t, 0, comp.Compare(
		f{Trips: 1, Samples: 100},
		f{Trips: 1, Samples: 100},
	))
	assert.Equal(t, 0, comp.Compare(
		f{Trips: 2, Samples: 100},
		f{Trips: 1, Samples: 200},
	))
	assert.Equal(t, 0, comp.Compare(
		f{Trips: 1, Samples: 200},
		f{Trips: 2, Samples: 100},
	))

	assert.Equal(t, 1, comp.Compare(
		f{Trips: 2, Samples: 100},
		f{Trips: 1, Samples: 100},
	))
	assert.Equal(t, 1, comp.Compare(
		f{Trips: 1, Samples: 200},
		f{Trips: 1, Samples: 100},
	))
}
