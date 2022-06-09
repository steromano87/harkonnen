package project

type Directory interface {
	Dir() string
	Create(workdir string) error
}
