package spec

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type EnsureSpec struct {
	Spec Specification
}

type RemoveSpec struct {
	Spec Specification
}

type ReplaceSpec struct {
	Spec Specification
}

// optional interface for specs that support removal
type RemovableSpec interface {
	Exists(project *Project) (bool, error)
	Remove(project *Project) error
}

// optional interface for specs that support replacement
type ReplaceableSpec interface {
	Equals(project *Project) (bool, error)
	Replace(project *Project) error
}

// PresentSpec methods
func (m *EnsureSpec) Check(project *Project) (bool, error) {
	return m.Spec.Check(project)
}

func (m *EnsureSpec) Apply(project *Project) error {
	return m.Spec.Apply(project)
}

// RemoveSpec methods
func (m *RemoveSpec) Check(project *Project) (bool, error) {
	if rm, ok := m.Spec.(RemovableSpec); ok {
		exists, err := rm.Exists(project)
		if err != nil {
			return false, err
		}
		return !exists, nil
	}

	log.Warn().Type("spec", m.Spec).Msg("Spec does not support removal; using fallback check")

	// fallback to inverted Check logic
	exists, err := m.Spec.Check(project)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (m *RemoveSpec) Apply(project *Project) error {
	log.Trace().Str("project", project.Name).Msg("Applying removal spec")

	if rm, ok := m.Spec.(RemovableSpec); ok {
		return rm.Remove(project)
	}

	log.Error().Type("spec", m.Spec).Msg("Spec does not support removal")
	return fmt.Errorf("spec type %T does not support removal", m.Spec)
}

// ReplaceSpec methods
func (m *ReplaceSpec) Check(project *Project) (bool, error) {
	if repl, ok := m.Spec.(ReplaceableSpec); ok {
		equal, err := repl.Equals(project)
		if err != nil {
			return false, err
		}
		return equal, nil
	}

	log.Warn().Type("spec", m.Spec).Msg("Spec does not support replacement; using fallback check")

	// fallback to standard Check logic
	return m.Spec.Check(project)
}

func (m *ReplaceSpec) Apply(project *Project) error {
	log.Trace().Str("project", project.Name).Msg("Applying replacement spec")

	if repl, ok := m.Spec.(ReplaceableSpec); ok {
		return repl.Replace(project)
	}

	log.Error().Type("spec", m.Spec).Msg("Spec does not support replacement")
	return fmt.Errorf("spec type %T does not support replacement", m.Spec)
}
