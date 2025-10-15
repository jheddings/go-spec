package spec

import (
	"testing"
)

func TestBasicProjectInfo(t *testing.T) {
	project := NewProject("test proj").
		WithDescription("test desc").
		WithPath("/path/to/proj").
		WithHomepage("https://test.com").
		WithVar("test-var", "test-val").
		Build()

	if project == nil {
		t.Fatal("project is nil")
	}

	if project.Name != "test proj" {
		t.Fatal("project name is not test proj")
	}

	if project.Desc != "test desc" {
		t.Fatal("project description is not test desc")
	}

	if project.Path != "/path/to/proj" {
		t.Fatal("project path is not /path/to/proj")
	}

	if project.URL != "https://test.com" {
		t.Fatal("project URL is not https://test.com")
	}

	if val, ok := project.Vars["test-var"]; !ok {
		t.Fatal("project var is not test-var")
	} else {
		if val != "test-val" {
			t.Fatal("project var is not test-val")
		}
	}
}
