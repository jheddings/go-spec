package gospec

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type SpecFactory func(config any) (Specification, error)

type Specification interface {
	Check(project *Project) (bool, error)
	Apply(project *Project) error
}

type SpecRegistry struct {
	lock  sync.RWMutex
	specs map[string]SpecFactory
}

var specRegistry = &SpecRegistry{
	specs: make(map[string]SpecFactory),
}

func RegisterSpec(name string, factory SpecFactory) {
	log.Trace().Str("spec", name).Msg("Registering specification")
	specRegistry.lock.Lock()
	defer specRegistry.lock.Unlock()
	specRegistry.specs[name] = factory
}

func CreateSpec(name string, config any) (Specification, error) {
	factory, ok := specRegistry.specs[name]
	if !ok {
		return nil, fmt.Errorf("specification %s not found", name)
	}
	return factory(config)
}
