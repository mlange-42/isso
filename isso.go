package isso

import (
	"log"
	"slices"
)

type subject int
type matrix int

// Requirement definition.
type Requirement struct {
	Subject string
	Matrix  string
	Samples int
	Times   []int
}

// Action definition.
type Action struct {
	Subject       string
	Matrix        string
	Reuse         string
	Time          int
	Samples       int
	TargetSamples int
}

// requirement for internal use, using no strings.
type requirement struct {
	Subject subject
	Matrix  matrix
	Samples int
	Times   []int
}

// action for internal use, using no strings.
type action struct {
	Subject       subject
	Matrix        matrix
	Reuse         subject
	IsReuse       bool
	Time          int
	Samples       int
	TargetSamples int
}

// Matrix definition.
type Matrix struct {
	Name     string
	CanReuse []string
}

// Actions of an internal solution.
type Actions struct {
	Actions []action
}

// Solution, translated back to using strings for subject and matrix.
type Solution[F any] struct {
	Actions []Action
	Fitness F
}

// Problem definition.
type Problem struct {
	subjectIDs   map[string]subject
	subjectNames map[subject]string
	matrixIDs    map[string]matrix
	matrixNames  map[matrix]string
	capacity     []int
	reusable     [][]bool
	requirements []requirement
}

// NewProblem creates a new problem definition.
func NewProblem(
	subjects []string,
	matrices []Matrix,
	capacity []int,
	requirements []Requirement) Problem {

	matrixIDs := map[string]matrix{}
	matrixNames := map[matrix]string{}
	for i, m := range matrices {
		matrixIDs[m.Name] = matrix(i)
		matrixNames[matrix(i)] = m.Name
	}

	subjectIDs := map[string]subject{}
	subjectNames := map[subject]string{}
	for i, s := range subjects {
		subjectIDs[s] = subject(i)
		subjectNames[subject(i)] = s
	}

	reusable := make([][]bool, len(matrices))
	for i, m := range matrices {
		reusable[i] = make([]bool, len(matrices))
		reusable[i][i] = true
		for _, ru := range m.CanReuse {
			if id, ok := matrixIDs[ru]; ok {
				reusable[i][id] = true
			} else {
				log.Fatalf("unknown matrix '%v'", ru)
			}
		}
	}

	req := make([]requirement, len(requirements))
	uniqueReq := map[subject]bool{}
	for i, r := range requirements {
		subject, ok := subjectIDs[r.Subject]
		if !ok {
			log.Fatalf("unknown subject '%v'", r.Subject)
		}

		if _, ok := uniqueReq[subject]; ok {
			log.Fatalf("duplicate subject '%v' in requirements", r.Subject)
		}
		uniqueReq[subject] = true

		matrix, ok := matrixIDs[r.Matrix]
		if !ok {
			log.Fatalf("unknown matrix '%v'", r.Matrix)
		}

		times := slices.Clone(r.Times)
		slices.Sort(times)
		times = slices.Compact(times)
		if len(times) != len(r.Times) {
			log.Fatalf("duplicate time entry in times for subject '%v'", r.Subject)
		}

		req[i] = requirement{
			Subject: subject,
			Matrix:  matrix,
			Samples: r.Samples,
			Times:   times,
		}
	}

	return Problem{
		subjectIDs:   subjectIDs,
		subjectNames: subjectNames,
		matrixIDs:    matrixIDs,
		matrixNames:  matrixNames,
		capacity:     capacity,
		reusable:     reusable,
		requirements: req,
	}
}

// Comparator interface or comparing fitness values.
type Comparator[F any] interface {
	Less(a, b F) bool
}

// Evaluator interface or deriving fitness from a solution.
type Evaluator[F any] interface {
	Evaluate(s *Actions) F
}

// Solver for optimization.
type Solver[F any] struct {
	problem       *Problem
	bestSolution  Actions
	bestFitness   F
	preserved     []action
	preservedTemp []action
	anySolution   bool
	evaluator     Evaluator[F]
	comparator    Comparator[F]
}

// NewSolver creates a new solver for a given fitness function.
func NewSolver[S comparable, M comparable, F any](evaluator Evaluator[F], comparator Comparator[F]) Solver[F] {
	return Solver[F]{
		evaluator:  evaluator,
		comparator: comparator,
	}
}

// Solve the given problem.
func (s *Solver[F]) Solve(problem *Problem) (Solution[F], bool) {
	s.problem = problem
	s.bestSolution = Actions{}
	s.preserved = []action{}
	s.preservedTemp = []action{}
	s.anySolution = false

	s.solve(&Actions{})

	if s.anySolution {
		actions := make([]Action, len(s.preserved))

		for i := range s.preserved {
			a := &s.preserved[i]
			var reuse string
			if a.IsReuse {
				reuse = s.problem.subjectNames[a.Reuse]
			}
			actions[i] = Action{
				Subject:       s.problem.subjectNames[a.Subject],
				Matrix:        s.problem.matrixNames[a.Matrix],
				Samples:       a.Samples,
				TargetSamples: a.TargetSamples,
				Time:          a.Time,
				Reuse:         reuse,
			}
		}

		return Solution[F]{
			Actions: actions,
			Fitness: s.bestFitness,
		}, true
	}
	return Solution[F]{}, false
}

func (s *Solver[F]) solve(sol *Actions) {
	fitness := s.evaluator.Evaluate(sol)
	if !s.comparator.Less(fitness, s.bestFitness) {
		return
	}

	var unsatisfied *requirement = nil
	var requiredSamples = 0

	capacity := slices.Clone(s.problem.capacity)

	for r := range s.problem.requirements {
		req := &s.problem.requirements[r]
		samples := req.Samples
		for a := range sol.Actions {
			act := &sol.Actions[a]

			if !slices.Contains(req.Times, act.Time) {
				continue
			}

			if !s.problem.reusable[req.Matrix][act.Matrix] {
				continue
			}

			equivalentSamples := min(act.Samples, samples)

			ownSample := req.Subject == act.Subject
			if ownSample {
				equivalentSamples = min(equivalentSamples, capacity[act.Time])
			}

			samples -= equivalentSamples

			if equivalentSamples > 0 {
				s.preservedTemp = append(s.preservedTemp, action{
					Subject:       req.Subject,
					Matrix:        req.Matrix,
					Samples:       equivalentSamples,
					Time:          act.Time,
					TargetSamples: req.Samples,
					IsReuse:       !ownSample,
					Reuse:         act.Subject,
				})
				if ownSample {
					capacity[act.Time] -= equivalentSamples
				}
			}
			if samples == 0 {
				break
			}
		}

		if samples > 0 {
			if unsatisfied == nil {
				unsatisfied = req
				requiredSamples = samples
			} else {
				if req.Matrix == unsatisfied.Matrix {
					if samples > requiredSamples {
						unsatisfied = req
						requiredSamples = samples
					}
				} else if s.problem.reusable[unsatisfied.Matrix][req.Matrix] {
					unsatisfied = req
					requiredSamples = samples
				}
			}
		}
	}

	if unsatisfied != nil {
		for _, t := range unsatisfied.Times {
			if capacity[t] <= 0 {
				continue
			}

			maxSamples := min(requiredSamples, capacity[t])

			sol.Actions = append(sol.Actions, action{
				Subject:       unsatisfied.Subject,
				Matrix:        unsatisfied.Matrix,
				Samples:       maxSamples,
				TargetSamples: unsatisfied.Samples,
				Time:          t,
				IsReuse:       false,
			})

			s.preservedTemp = s.preservedTemp[:0]
			s.solve(sol)

			sol.Actions = sol.Actions[:len(sol.Actions)-1]
		}
	} else {
		if s.comparator.Less(fitness, s.bestFitness) {
			s.bestFitness = fitness
			s.bestSolution = Actions{
				Actions: slices.Clone(sol.Actions),
			}
			s.preserved = slices.Clone(s.preservedTemp)
			s.preservedTemp = s.preservedTemp[:0]
			s.anySolution = true
		}
	}
}
