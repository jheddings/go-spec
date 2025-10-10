package spec

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

type DeferredSpec struct {
	SpecFunc func() Specification

	lock sync.Mutex
	spec *Specification
}

func (d *DeferredSpec) init() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.spec = Ptr(d.SpecFunc())
}

func (d *DeferredSpec) Check(project *Project) (bool, error) {
	d.init()

	if d.spec == nil {
		return false, fmt.Errorf("unable to initialize deferred spec")
	}

	return (*d.spec).Check(project)
}

func (d *DeferredSpec) Apply(project *Project) error {
	d.init()

	if d.spec == nil {
		return fmt.Errorf("unable to initialize deferred spec")
	}

	return (*d.spec).Apply(project)
}
