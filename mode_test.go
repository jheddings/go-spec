package spec

import "testing"

type TestModeSpec struct {
	apply   bool
	check   bool
	equals  bool
	exists  bool
	replace bool
	remove  bool
}

func (t *TestModeSpec) Check(project *Project) (bool, error) {
	t.check = true
	return false, nil
}

func (t *TestModeSpec) Apply(project *Project) error {
	t.apply = true
	return nil
}

func (t *TestModeSpec) Equals(project *Project) (bool, error) {
	t.equals = true
	return false, nil
}

func (t *TestModeSpec) Exists(project *Project) (bool, error) {
	t.exists = true
	return true, nil
}

func (t *TestModeSpec) Replace(project *Project) error {
	t.replace = true
	return nil
}

func (t *TestModeSpec) Remove(project *Project) error {
	t.remove = true
	return nil
}

func TestEnsureSpec(t *testing.T) {
	spec := &TestModeSpec{}
	enSpec := &EnsureSpec{Spec: spec}

	project := NewProject("test-ensure").WithSpec(enSpec).Build()
	project.BuildAll()

	if !spec.check {
		t.Fatal("failed to check")
	}

	if !spec.apply {
		t.Fatal("failed to apply")
	}
}

func TestRemoveSpec(t *testing.T) {
	spec := &TestModeSpec{}
	rmSpec := &RemoveSpec{Spec: spec}

	project := NewProject("test-remove").WithSpec(rmSpec).Build()
	project.BuildAll()

	if spec.check {
		t.Fatal("should not have checked")
	}

	if !spec.exists {
		t.Fatal("failed to check exists")
	}

	if !spec.remove {
		t.Fatal("failed to remove")
	}
}

func TestReplaceSpec(t *testing.T) {
	spec := &TestModeSpec{}
	repSpec := &ReplaceSpec{Spec: spec}

	project := NewProject("test-replace").WithSpec(repSpec).Build()
	project.BuildAll()

	if spec.check {
		t.Fatal("should not have checked")
	}

	if !spec.equals {
		t.Fatal("failed to check equals")
	}

	if !spec.replace {
		t.Fatal("failed to replace")
	}
}
