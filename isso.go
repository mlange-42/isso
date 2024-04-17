package isso

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
	From          Subject
	Matrix        Matrix
	Time          int
	Samples       int
	TargetSamples int
}

type Problem struct {
	matrixIDs   map[string]Matrix
	matrixNames map[Matrix]string
	capacity    []int
}

func NewProblem(matrices []string, capacity []int) Problem {
	ids := map[string]Matrix{}
	names := map[Matrix]string{}
	for i, m := range matrices {
		ids[m] = Matrix(i)
		names[Matrix(i)] = m
	}
	return Problem{
		matrixIDs:   ids,
		matrixNames: names,
		capacity:    capacity,
	}
}
