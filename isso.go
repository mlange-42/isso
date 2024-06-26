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
	Times   []int
	Samples int
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
	Times   []int
	Subject subject
	Matrix  matrix
	Samples int
}

// ActionDef for internal use, using no strings.
// It needs to be public as it is used in fitness evaluators.
type ActionDef struct {
	Subject       subject
	Matrix        matrix
	Reuse         subject
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
type actions struct {
	Actions []ActionDef
}

// Solution, translated back to using strings for subject and matrix.
type Solution[F any] struct {
	Fitness F
	Actions []Action
}

// solution for internal use.
type solution[F any] struct {
	Fitness F
	Actions []ActionDef
}

type ProblemDef struct {
	Matrices     []Matrix
	Capacity     []int
	Requirements []Requirement
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
func NewProblem(problem ProblemDef) Problem {

	matrixIDs := map[string]matrix{}
	matrixNames := map[matrix]string{}
	for i, m := range problem.Matrices {
		matrixIDs[m.Name] = matrix(i)
		matrixNames[matrix(i)] = m.Name
	}

	reusable := make([][]bool, len(problem.Matrices))
	for i, m := range problem.Matrices {
		reusable[i] = make([]bool, len(problem.Matrices))
		reusable[i][i] = true
		for _, ru := range m.CanReuse {
			if id, ok := matrixIDs[ru]; ok {
				reusable[i][id] = true
			} else {
				log.Fatalf("unknown matrix '%v'", ru)
			}
		}
	}

	req := make([]requirement, len(problem.Requirements))
	subjectIDs := map[string]subject{}
	subjectNames := map[subject]string{}
	for i, r := range problem.Requirements {
		if _, ok := subjectIDs[r.Subject]; ok {
			log.Fatalf("duplicate subject '%v' in requirements", r.Subject)
		}
		sub := subject(i)

		subjectIDs[r.Subject] = sub
		subjectNames[sub] = r.Subject

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
			Subject: sub,
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
		capacity:     problem.Capacity,
		reusable:     reusable,
		requirements: req,
	}
}

// Comparator interface or comparing fitness values.
type Comparator[F any] interface {
	Compare(a, b F) int
	IsPareto() bool
}

// Evaluator interface or deriving fitness from a solution.
type Evaluator[F any] interface {
	Evaluate(s []ActionDef) F
}

// Solver for optimization.
type Solver[F comparable] struct {
	bestFitness  F
	evaluator    Evaluator[F]
	comparator   Comparator[F]
	problem      *Problem
	solutions    []solution[F]
	tempSolution []ActionDef
}

// NewSolver creates a new solver for a given fitness function.
func NewSolver[F comparable](evaluator Evaluator[F], comparator Comparator[F]) Solver[F] {
	return Solver[F]{
		evaluator:  evaluator,
		comparator: comparator,
	}
}

// Solve the given problem.
func (s *Solver[F]) Solve(problem *Problem) ([]Solution[F], bool) {
	s.problem = problem
	s.solutions = []solution[F]{}
	s.tempSolution = []ActionDef{}

	s.solve(&actions{})

	if len(s.solutions) > 0 {
		return s.toSolutions(), true
	}
	return []Solution[F]{}, false
}

// toSolutions converts the solution results to the solution output type,
// translating integer IDs back to strings.
func (s *Solver[F]) toSolutions() []Solution[F] {
	solutions := []Solution[F]{}

	for _, sol := range s.solutions {
		actions := make([]Action, len(sol.Actions))

		for i := range sol.Actions {
			a := &sol.Actions[i]
			var reuse string
			if a.Reuse >= 0 {
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

		solutions = append(solutions, Solution[F]{
			Actions: actions,
			Fitness: sol.Fitness,
		})
	}

	return solutions
}

// Recursive solver function.
func (s *Solver[F]) solve(sol *actions) {
	fitness := s.evaluator.Evaluate(sol.Actions)

	if s.comparator.IsPareto() {
		if !s.isParetoOptimal(fitness, false) {
			return
		}
	} else {
		if s.comparator.Compare(fitness, s.bestFitness) > 0 {
			return
		}
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
				reuse := subject(-1)
				if !ownSample {
					reuse = act.Subject
				}
				s.tempSolution = append(s.tempSolution, ActionDef{
					Subject:       req.Subject,
					Matrix:        req.Matrix,
					Samples:       equivalentSamples,
					Time:          act.Time,
					TargetSamples: req.Samples,
					Reuse:         reuse,
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
				// for the same matrix, prefer the larger sample.
				if req.Matrix == unsatisfied.Matrix {
					if samples > requiredSamples {
						unsatisfied = req
						requiredSamples = samples
					}
					// if not the same matrix, prefer the one that can be re-used by the other.
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

			sol.Actions = append(sol.Actions, ActionDef{
				Subject:       unsatisfied.Subject,
				Matrix:        unsatisfied.Matrix,
				Samples:       min(requiredSamples, capacity[t]),
				TargetSamples: unsatisfied.Samples,
				Time:          t,
				Reuse:         -1,
			})
			s.tempSolution = s.tempSolution[:0]

			s.solve(sol)

			sol.Actions = sol.Actions[:len(sol.Actions)-1]
		}
	} else {
		if s.comparator.IsPareto() {
			if s.isParetoOptimal(fitness, true) {
				s.solutions = append(s.solutions, solution[F]{
					Actions: slices.Clone(s.tempSolution),
					Fitness: fitness,
				})
				s.tempSolution = s.tempSolution[:0]
			}
		} else {
			comp := s.comparator.Compare(fitness, s.bestFitness)
			if comp < 0 {
				s.solutions = s.solutions[:0]
			}
			if comp <= 0 {
				s.bestFitness = fitness
				s.solutions = append(s.solutions, solution[F]{
					Actions: slices.Clone(s.tempSolution),
					Fitness: fitness,
				})
				s.tempSolution = s.tempSolution[:0]
			}
		}
	}
}

// removeSolution swap-removes the solution at the given index.
func (s *Solver[F]) removeSolution(idx int) {
	ln := len(s.solutions) - 1
	s.solutions[idx], s.solutions[ln] = s.solutions[ln], s.solutions[idx]
	s.solutions = s.solutions[:ln]
}

// isParetoOptimal checks if the given fitness is pareto-optimal and should be retained as a solution.
// If argument remove is true, all non-optimal solutions are removed from the Solver.
func (s *Solver[F]) isParetoOptimal(f F, remove bool) bool {
	betterThanAny := len(s.solutions) == 0
	hasDuplicate := false
	for i := len(s.solutions) - 1; i >= 0; i-- {
		comp := s.comparator.Compare(f, s.solutions[i].Fitness)

		if comp < 0 {
			if remove {
				s.removeSolution(i)
			}
		} else if f == s.solutions[i].Fitness {
			hasDuplicate = true
		}

		if comp <= 0 {
			betterThanAny = true
		} else {
			hasDuplicate = true
		}
	}
	return betterThanAny && !hasDuplicate
}
