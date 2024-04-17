package isso

import (
	"log"
	"slices"
)

type Subject int
type Matrix int

type RequirementDef[S comparable, M comparable] struct {
	Subject S
	Matrix  M
	Samples int
	Times   []int
}

type Requirement struct {
	Subject Subject
	Matrix  Matrix
	Samples int
	Times   []int
}

type Action struct {
	Subject       Subject
	Matrix        Matrix
	Reuse         Subject
	IsReuse       bool
	Time          int
	Samples       int
	TargetSamples int
}

type MatrixDef[M comparable] struct {
	Name     M
	CanReuse []M
}

type Solution struct {
	Actions []Action
}

type Problem[S comparable, M comparable] struct {
	subjectIDs   map[S]Subject
	subjectNames map[Subject]S
	matrixIDs    map[M]Matrix
	matrixNames  map[Matrix]M
	capacity     []int
	reusable     [][]bool
	requirements []Requirement
}

func NewProblem[S comparable, M comparable](
	subjects []S,
	matrices []MatrixDef[M],
	capacity []int,
	requirements []RequirementDef[S, M]) Problem[S, M] {

	matrixIDs := map[M]Matrix{}
	matrixNames := map[Matrix]M{}
	for i, m := range matrices {
		matrixIDs[m.Name] = Matrix(i)
		matrixNames[Matrix(i)] = m.Name
	}

	subjectIDs := map[S]Subject{}
	subjectNames := map[Subject]S{}
	for i, s := range subjects {
		subjectIDs[s] = Subject(i)
		subjectNames[Subject(i)] = s
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

	req := make([]Requirement, len(requirements))
	uniqueReq := map[Subject]bool{}
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

		req[i] = Requirement{
			Subject: subject,
			Matrix:  matrix,
			Samples: r.Samples,
			Times:   times,
		}
	}

	return Problem[S, M]{
		subjectIDs:   subjectIDs,
		subjectNames: subjectNames,
		matrixIDs:    matrixIDs,
		matrixNames:  matrixNames,
		capacity:     capacity,
		reusable:     reusable,
		requirements: req,
	}
}

type Comparator[F any] interface {
	Less(a, b F) bool
}
type Evaluator[F any] interface {
	Evaluate(s *Solution) F
}

type Solver[S comparable, M comparable, F any] struct {
	problem       *Problem[S, M]
	bestSolution  []Action
	preserved     []Action
	preservedTemp []Action
	anySolution   bool
	evaluator     Evaluator[F]
	comparator    Comparator[F]
}

func NewSolver[S comparable, M comparable, F any](evaluator Evaluator[F], comparator Comparator[F]) Solver[S, M, F] {
	return Solver[S, M, F]{
		evaluator:  evaluator,
		comparator: comparator,
	}
}

func (s *Solver[S, M, F]) Solve(problem *Problem[S, M]) {
	s.problem = problem
	s.bestSolution = []Action{}
	s.preserved = []Action{}
	s.preservedTemp = []Action{}
	s.anySolution = false

	s.solve(&Solution{})
}

func (s *Solver[S, M, F]) solve(solution *Solution) {

}
