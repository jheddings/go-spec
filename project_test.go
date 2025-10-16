package spec

import (
	"fmt"
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

func TestBuildAll(t *testing.T) {
	t.Run("build all with single spec", func(t *testing.T) {
		spec := &TestSpec{}
		project := NewProject("test").
			WithSpec(spec).
			Build()

		err := project.BuildAll()
		if err != nil {
			t.Fatalf("BuildAll failed: %v", err)
		}

		if !spec.check {
			t.Fatal("spec Check was not called")
		}

		if !spec.apply {
			t.Fatal("spec Apply was not called")
		}
	})

	t.Run("build all with multiple specs", func(t *testing.T) {
		spec1 := &TestSpec{}
		spec2 := &TestSpec{}

		project := NewProject("test").
			WithSpec(spec1).
			WithSpec(spec2).
			Build()

		err := project.BuildAll()
		if err != nil {
			t.Fatalf("BuildAll failed: %v", err)
		}

		if !spec1.apply || !spec2.apply {
			t.Fatal("not all specs were applied")
		}
	})

	t.Run("build all skips when check returns true", func(t *testing.T) {
		spec := &TestSpec{check: true}
		project := NewProject("test").
			WithSpec(spec).
			Build()

		err := project.BuildAll()
		if err != nil {
			t.Fatalf("BuildAll failed: %v", err)
		}

		if !spec.check {
			t.Fatal("failed to check")
		}

		if spec.apply {
			t.Fatal("should not have applied")
		}
	})

	t.Run("build all stops on check error", func(t *testing.T) {
		spec := &TestCheckErrorSpec{}
		project := NewProject("test").
			WithSpec(spec).
			Build()

		err := project.BuildAll()
		if err == nil {
			t.Fatal("expected error from BuildAll")
		}

		expectedErr := "check error"
		if err.Error() != expectedErr {
			t.Fatalf("expected error %q, got %q", expectedErr, err.Error())
		}
	})

	t.Run("build all stops on apply error", func(t *testing.T) {
		spec := &TestApplyErrorSpec{}
		project := NewProject("test").
			WithSpec(spec).
			Build()

		err := project.BuildAll()
		if err == nil {
			t.Fatal("expected error from BuildAll")
		}

		expectedErr := "apply error"
		if err.Error() != expectedErr {
			t.Fatalf("expected error %q, got %q", expectedErr, err.Error())
		}
	})
}

// Test helpers

type TestCheckErrorSpec struct{}

func (t *TestCheckErrorSpec) Check(project *Project) (bool, error) {
	return false, fmt.Errorf("check error")
}

func (t *TestCheckErrorSpec) Apply(project *Project) error {
	return nil
}

type TestApplyErrorSpec struct{}

func (t *TestApplyErrorSpec) Check(project *Project) (bool, error) {
	return false, nil
}

func (t *TestApplyErrorSpec) Apply(project *Project) error {
	return fmt.Errorf("apply error")
}
