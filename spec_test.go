package spec

import (
	"fmt"
	"testing"
)

func TestRegisterSpec(t *testing.T) {
	testFactory := func(config any) (Specification, error) {
		return &TestSpec{}, nil
	}

	RegisterSpec("test-spec", testFactory)

	specRegistry.lock.RLock()
	defer specRegistry.lock.RUnlock()

	factory, ok := specRegistry.specs["test-spec"]
	if !ok {
		t.Fatal("spec not found in registry")
	}

	if factory == nil {
		t.Fatal("factory is nil")
	}

	// clean up - need write lock
	specRegistry.lock.RUnlock()
	specRegistry.lock.Lock()
	delete(specRegistry.specs, "test-spec")
	specRegistry.lock.Unlock()
	specRegistry.lock.RLock()
}

func TestCreateSpec(t *testing.T) {
	t.Run("create existing spec", func(t *testing.T) {
		testFactory := func(config any) (Specification, error) {
			return &TestSpec{}, nil
		}
		RegisterSpec("create-test-spec", testFactory)

		spec, err := CreateSpec("create-test-spec", nil)
		if err != nil {
			t.Fatalf("CreateSpec failed: %v", err)
		}

		if spec == nil {
			t.Fatal("CreateSpec returned nil spec")
		}

		if _, ok := spec.(*TestSpec); !ok {
			t.Fatal("expected TestModeSpec")
		}

		// clean up
		specRegistry.lock.Lock()
		delete(specRegistry.specs, "create-test-spec")
		specRegistry.lock.Unlock()
	})

	t.Run("create non-existent spec", func(t *testing.T) {
		spec, err := CreateSpec("non-existent-spec", nil)
		if err == nil {
			t.Fatal("expected error for non-existent spec")
		}

		if spec != nil {
			t.Fatal("expected nil spec")
		}

		expectedErr := "specification non-existent-spec not found"
		if err.Error() != expectedErr {
			t.Fatalf("expected error %q, got %q", expectedErr, err.Error())
		}
	})

	t.Run("factory returns error", func(t *testing.T) {
		errorFactory := func(config any) (Specification, error) {
			return nil, fmt.Errorf("factory error")
		}
		RegisterSpec("error-spec", errorFactory)

		spec, err := CreateSpec("error-spec", nil)
		if err == nil {
			t.Fatal("expected error from factory")
		}

		if spec != nil {
			t.Fatal("expected nil spec")
		}

		// clean up
		specRegistry.lock.Lock()
		delete(specRegistry.specs, "error-spec")
		specRegistry.lock.Unlock()
	})
}

func TestDeferredSpec(t *testing.T) {
	t.Run("deferred spec check", func(t *testing.T) {
		testSpec := &TestSpec{}
		deferred := &DeferredSpec{
			SpecFunc: func() Specification {
				return testSpec
			},
		}

		project := NewProject("test").Build()

		result, err := deferred.Check(project)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}

		if result {
			t.Fatal("expected Check to return false")
		}

		if !testSpec.check {
			t.Fatal("underlying spec Check was not called")
		}
	})

	t.Run("deferred spec apply", func(t *testing.T) {
		testSpec := &TestSpec{}
		deferred := &DeferredSpec{
			SpecFunc: func() Specification {
				return testSpec
			},
		}

		project := NewProject("test").Build()

		err := deferred.Apply(project)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		if !testSpec.apply {
			t.Fatal("underlying spec Apply was not called")
		}
	})
}
