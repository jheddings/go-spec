package spec

import (
	"slices"

	"github.com/rs/zerolog/log"
)

type Project struct {
	Name  string
	Desc  string
	Path  string
	URL   string
	Vars  map[string]any
	Specs []Specification
}

type ProjectBuilder struct {
	project *Project
}

var projectsConfig = []Project{}

func NewProject(name string) *ProjectBuilder {
	return &ProjectBuilder{
		project: &Project{
			Name:  name,
			Desc:  "",
			Vars:  make(map[string]any),
			Specs: []Specification{},
		},
	}
}

func (p *ProjectBuilder) WithDescription(description string) *ProjectBuilder {
	p.project.Desc = description
	return p
}

func (p *ProjectBuilder) WithPath(path string) *ProjectBuilder {
	p.project.Path = path
	return p
}

func (p *ProjectBuilder) WithHomepage(url string) *ProjectBuilder {
	p.project.URL = url
	return p
}

func (p *ProjectBuilder) WithVar(name string, value any) *ProjectBuilder {
	p.project.Vars[name] = value
	return p
}

func (p *ProjectBuilder) WithSpec(spec Specification) *ProjectBuilder {
	p.project.Specs = append(p.project.Specs, spec)
	return p
}

func (b *ProjectBuilder) WithSpecPresent(spec Specification) *ProjectBuilder {
	return b.WithSpec(&EnsureSpec{Spec: spec})
}

func (b *ProjectBuilder) WithSpecRemove(spec Specification) *ProjectBuilder {
	return b.WithSpec(&RemoveSpec{Spec: spec})
}

func (b *ProjectBuilder) WithSpecReplace(spec Specification) *ProjectBuilder {
	return b.WithSpec(&ReplaceSpec{Spec: spec})
}

func (p *ProjectBuilder) WithBlueprint(bp Blueprint) *ProjectBuilder {
	p.project.Specs = append(p.project.Specs, bp.Specs...)
	return p
}

func (p *ProjectBuilder) Build() *Project {
	return p.project
}

func (p *Project) BuildAll() error {
	for _, spec := range p.Specs {
		if err := p.applySpec(spec); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) applySpec(spec Specification) error {
	check, err := spec.Check(p)
	if err != nil {
		log.Warn().Str("project", p.Name).Type("spec", spec).Msg("Failed to check")
		return err
	}

	if check {
		log.Info().Str("project", p.Name).Type("spec", spec).Msg("Skipping; up to date")
		return nil
	}

	log.Info().Str("project", p.Name).Type("spec", spec).Msg("Applying")
	if err := spec.Apply(p); err != nil {
		log.Warn().Str("project", p.Name).Type("spec", spec).Msg("Failed to apply")
		return err
	}

	return nil
}

func RegisterProject(project *Project) {
	log.Trace().Str("project", project.Name).Msg("Registering project")
	projectsConfig = append(projectsConfig, *project)
}

func FilterProjects(names []string) []*Project {
	projects := []*Project{}
	for _, project := range projectsConfig {
		if slices.Contains(names, project.Name) || len(names) == 0 {
			projects = append(projects, &project)
		}
	}
	return projects
}
