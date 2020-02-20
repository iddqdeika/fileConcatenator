package definition

type Validator func(path string) (ok bool)

type Handler func(path string)

type Processor interface {
	WithHandler(validator Validator, handler Handler) Processor
	ProcessPath(path string) error
}
