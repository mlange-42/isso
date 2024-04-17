package isso

import (
	"log"
)

type Subject int
type Matrix int

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
	Time          int
	Samples       int
	TargetSamples int
}

type MatrixDef[M comparable] struct {
	Name     M
	CanReuse []M
}

type Problem[M comparable] struct {
	matrixIDs   map[M]Matrix
	matrixNames map[Matrix]M
	capacity    []int
	reusable    [][]bool
}

func NewProblem[M comparable](matrices []MatrixDef[M], capacity []int) Problem[M] {
	ids := map[M]Matrix{}
	names := map[Matrix]M{}
	for i, m := range matrices {
		ids[m.Name] = Matrix(i)
		names[Matrix(i)] = m.Name
	}

	reusable := make([][]bool, len(matrices))
	for i, m := range matrices {
		reusable[i] = make([]bool, len(matrices))
		reusable[i][i] = true
		for _, ru := range m.CanReuse {
			if id, ok := ids[ru]; ok {
				reusable[i][id] = true
			} else {
				log.Fatalf("unknown matrix '%v'", ru)
			}
		}
	}

	return Problem[M]{
		matrixIDs:   ids,
		matrixNames: names,
		capacity:    capacity,
		reusable:    reusable,
	}
}
