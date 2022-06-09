package project

type File interface {
	File() string
	Create(workdir string) error
}
