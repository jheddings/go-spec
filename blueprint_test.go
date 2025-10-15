package spec

import "testing"

func TestRegisterBlueprint(t *testing.T) {
	bp := &Blueprint{
		Name:  "test-blueprint",
		Specs: []Specification{},
	}

	RegisterBlueprint(bp)

	blueprintRegistry.lock.RLock()
	defer blueprintRegistry.lock.RUnlock()

	registered, ok := blueprintRegistry.blueprints["test-blueprint"]
	if !ok {
		t.Fatal("blueprint not found in registry")
	}

	if registered.Name != bp.Name {
		t.Fatalf("expected blueprint name %s, got %s", bp.Name, registered.Name)
	}

	// clean up
	blueprintRegistry.lock.RUnlock()
	blueprintRegistry.lock.Lock()
	delete(blueprintRegistry.blueprints, "test-blueprint")
	blueprintRegistry.lock.Unlock()
	blueprintRegistry.lock.RLock()
}

func TestBlueprintBuilder(t *testing.T) {
	t.Run("basic blueprint", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 0 {
			t.Fatalf("expected 0 specs, got %d", len(bp.Specs))
		}

		if bp.Name != "test-bp" {
			t.Fatalf("expected name test-bp, got %s", bp.Name)
		}
	})

	t.Run("with single spec", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.WithSpec(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 1 {
			t.Fatalf("expected 1 specs, got %d", len(bp.Specs))
		}

		if _, ok := bp.Specs[0].(*TestSpec); !ok {
			t.Fatal("expected TestModeSpec")
		}
	})

	t.Run("with multiple specs", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.
			WithSpec(&TestSpec{}).
			WithSpec(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 2 {
			t.Fatalf("expected 2 specs, got %d", len(bp.Specs))
		}
	})

	t.Run("with deferred spec", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.WithDeferredSpec(func() Specification {
			return &TestSpec{}
		})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 1 {
			t.Fatalf("expected 1 specs, got %d", len(bp.Specs))
		}

		if _, ok := bp.Specs[0].(*DeferredSpec); !ok {
			t.Fatal("expected DeferredSpec")
		}
	})

	t.Run("with spec present", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.WithSpecPresent(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 1 {
			t.Fatalf("expected 1 specs, got %d", len(bp.Specs))
		}

		if _, ok := bp.Specs[0].(*EnsureSpec); !ok {
			t.Fatal("expected EnsureSpec")
		}
	})

	t.Run("with spec remove", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.WithSpecRemove(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 1 {
			t.Fatalf("expected 1 specs, got %d", len(bp.Specs))
		}

		if _, ok := bp.Specs[0].(*RemoveSpec); !ok {
			t.Fatal("expected RemoveSpec")
		}
	})

	t.Run("with spec replace", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		builder = builder.WithSpecReplace(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 1 {
			t.Fatalf("expected 1 specs, got %d", len(bp.Specs))
		}

		if _, ok := bp.Specs[0].(*ReplaceSpec); !ok {
			t.Fatal("expected ReplaceSpec")
		}
	})

	t.Run("with nested blueprint", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		nested := Blueprint{
			Name: "nested-bp",
			Specs: []Specification{
				&TestSpec{},
				&TestSpec{},
			},
		}
		builder = builder.WithBlueprint(nested)
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 2 {
			t.Fatalf("expected 2 specs, got %d", len(bp.Specs))
		}
	})

	t.Run("complex blueprint", func(t *testing.T) {
		builder := NewBlueprint("test-bp")
		nested := Blueprint{
			Name: "nested",
			Specs: []Specification{
				&TestSpec{},
			},
		}
		builder = builder.
			WithSpec(&TestSpec{}).
			WithSpecPresent(&TestSpec{}).
			WithBlueprint(nested).
			WithSpecRemove(&TestSpec{})
		bp := builder.Build()

		if bp == nil {
			t.Fatal("Build returned nil")
		}

		if len(bp.Specs) != 4 {
			t.Fatalf("expected 4 specs, got %d", len(bp.Specs))
		}
	})
}
