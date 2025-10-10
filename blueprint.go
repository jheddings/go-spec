package spec

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Blueprint struct {
	Name  string
	Specs []Specification
}

type BlueprintRegistry struct {
	lock       sync.RWMutex
	blueprints map[string]*Blueprint
}

type BlueprintBuilder struct {
	blueprint *Blueprint
}

var blueprintRegistry = &BlueprintRegistry{
	blueprints: make(map[string]*Blueprint),
}

func RegisterBlueprint(bp *Blueprint) {
	log.Trace().Str("blueprint", bp.Name).Msg("Registering blueprint")
	blueprintRegistry.lock.Lock()
	defer blueprintRegistry.lock.Unlock()
	blueprintRegistry.blueprints[bp.Name] = bp
}

func NewBlueprint(name string) *BlueprintBuilder {
	return &BlueprintBuilder{
		blueprint: &Blueprint{
			Name:  name,
			Specs: []Specification{},
		},
	}
}

func (b *BlueprintBuilder) WithSpec(spec Specification) *BlueprintBuilder {
	b.blueprint.Specs = append(b.blueprint.Specs, spec)
	return b
}

func (b *BlueprintBuilder) WithDeferredSpec(fn func() Specification) *BlueprintBuilder {
	return b.WithSpec(&DeferredSpec{SpecFunc: fn})
}

func (b *BlueprintBuilder) WithSpecPresent(spec Specification) *BlueprintBuilder {
	return b.WithSpec(&EnsureSpec{Spec: spec})
}

func (b *BlueprintBuilder) WithSpecRemove(spec Specification) *BlueprintBuilder {
	return b.WithSpec(&RemoveSpec{Spec: spec})
}

func (b *BlueprintBuilder) WithSpecReplace(spec Specification) *BlueprintBuilder {
	return b.WithSpec(&ReplaceSpec{Spec: spec})
}

func (b *BlueprintBuilder) WithBlueprint(bp Blueprint) *BlueprintBuilder {
	b.blueprint.Specs = append(b.blueprint.Specs, bp.Specs...)
	return b
}

func (b *BlueprintBuilder) Build() *Blueprint {
	return b.blueprint
}
