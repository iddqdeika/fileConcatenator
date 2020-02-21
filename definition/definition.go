package definition

type Validator func(path string, isDir bool) (ok bool)

type Handler func(path string)

type Processor interface {
	WithHandler(validator Validator, handler Handler) Processor
	ProcessPath(path string) error
}
