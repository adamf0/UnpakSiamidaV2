package domain

type IndikatorTree struct {
	IndikatorId       int
	ParentIndikatorId *int
	Pointing          string
}