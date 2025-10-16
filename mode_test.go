package spec

import "testing"

type TestSpec struct {
	apply   bool
	check   bool
	equals  bool
	exists  bool
	replace bool
	remove  bool
}

func (t *TestSpec) Check(project *Project) (bool, error) {
	orig := t.check
	t.check = true
	return orig, nil
}

func (t *TestSpec) Apply(project *Project) error {
	t.apply = true
	return nil
}

func (t *TestSpec) Equals(project *Project) (bool, error) {
	orig := t.equals
	t.equals = true
	return orig, nil
}

func (t *TestSpec) Exists(project *Project) (bool, error) {
	orig := t.exists
	t.exists = true
	return orig, nil
}

func (t *TestSpec) Replace(project *Project) error {
	t.replace = true
	return nil
}

func (t *TestSpec) Remove(project *Project) error {
	t.remove = true
	return nil
}

func TestModalSpecs(t *testing.T) {
	t.Run("ensure spec", func(t *testing.T) {
		spec := &TestSpec{}
		enSpec := &EnsureSpec{Spec: spec}

		project := NewProject("test-ensure").WithSpec(enSpec).Build()
		project.BuildAll()

		if !spec.check {
			t.Fatal("failed to check")
		}

		if !spec.apply {
			t.Fatal("failed to apply")
		}
	})

	t.Run("remove spec", func(t *testing.T) {
		spec := &TestSpec{exists: true}
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

	})

	t.Run("replace spec", func(t *testing.T) {
		spec := &TestSpec{}
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
	})
}
