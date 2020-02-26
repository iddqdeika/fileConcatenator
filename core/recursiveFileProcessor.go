package core

import (
	"fileConcatenator/definition"
	"fmt"
	"os"
)

type dispatcher struct {
	validate definition.Validator
	handle   definition.Handler
}

func NewRecursiveProcessor() definition.Processor {
	return &recursiveProcessor{}
}

type recursiveProcessor struct {
	dispatchers []*dispatcher
}

func (p *recursiveProcessor) WithHandler(validator definition.Validator, handler definition.Handler) definition.Processor {
	if validator == nil || handler == nil {
		return p
	}
	p.dispatchers = append(p.dispatchers, &dispatcher{
		validate: validator,
		handle:   handler,
	})
	return p
}

func (p *recursiveProcessor) ProcessPath(path string) error {
	return p.process(path)
}

func (p *recursiveProcessor) process(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("(path=%v) cant open file, err: %v", path, err)
	}
	s, err := f.Stat()
	if err != nil {
		return fmt.Errorf("(path=%v) cant get file/dir stat, err: %v", path, err)
	}
	if s.IsDir() {
		names, err := f.Readdirnames(0)
		if err != nil {
			return fmt.Errorf("(path=%v) cant get dirnames, err: %v", path, err)
		}
		for _, name := range names {
			p.process(path + string(os.PathSeparator) + name)
		}
	}
	p.dispatch(path, s.IsDir())
	return nil
}

func (p *recursiveProcessor) dispatch(path string, isDir bool) {
	for _, d := range p.dispatchers {
		if d.validate(path, isDir) {
			d.handle(path)
		}
	}
}
