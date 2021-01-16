package models

type IFilter interface {
	Filter(exercises []Exercise) []Exercise
}
