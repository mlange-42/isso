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

type MatrixDef struct {
	Name     string
	CanReuse []string
}

type Problem struct {
	matrixIDs   map[string]Matrix
	matrixNames map[Matrix]string
	capacity    []int
	reusable    [][]bool
}

func NewProblem(matrices []MatrixDef, capacity []int) Problem {
	ids := map[string]Matrix{}
	names := map[Matrix]string{}
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
				log.Fatalf("unknown matrix '%s'", ru)
			}
		}
	}

	return Problem{
		matrixIDs:   ids,
		matrixNames: names,
		capacity:    capacity,
		reusable:    reusable,
	}
}
