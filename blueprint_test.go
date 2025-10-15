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
	tests := []struct {
		name          string
		buildFunc     func(*BlueprintBuilder) *BlueprintBuilder
		expectedSpecs int
		validate      func(*testing.T, *Blueprint)
	}{
		{
			name: "basic blueprint",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b
			},
			expectedSpecs: 0,
			validate: func(t *testing.T, bp *Blueprint) {
				if bp.Name != "test-bp" {
					t.Fatalf("expected name test-bp, got %s", bp.Name)
				}
			},
		},
		{
			name: "with single spec",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.WithSpec(&TestSpec{})
			},
			expectedSpecs: 1,
			validate: func(t *testing.T, bp *Blueprint) {
				if _, ok := bp.Specs[0].(*TestSpec); !ok {
					t.Fatal("expected TestModeSpec")
				}
			},
		},
		{
			name: "with multiple specs",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.
					WithSpec(&TestSpec{}).
					WithSpec(&TestSpec{})
			},
			expectedSpecs: 2,
			validate:      nil,
		},
		{
			name: "with deferred spec",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.WithDeferredSpec(func() Specification {
					return &TestSpec{}
				})
			},
			expectedSpecs: 1,
			validate: func(t *testing.T, bp *Blueprint) {
				if _, ok := bp.Specs[0].(*DeferredSpec); !ok {
					t.Fatal("expected DeferredSpec")
				}
			},
		},
		{
			name: "with spec present",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.WithSpecPresent(&TestSpec{})
			},
			expectedSpecs: 1,
			validate: func(t *testing.T, bp *Blueprint) {
				if _, ok := bp.Specs[0].(*EnsureSpec); !ok {
					t.Fatal("expected EnsureSpec")
				}
			},
		},
		{
			name: "with spec remove",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.WithSpecRemove(&TestSpec{})
			},
			expectedSpecs: 1,
			validate: func(t *testing.T, bp *Blueprint) {
				if _, ok := bp.Specs[0].(*RemoveSpec); !ok {
					t.Fatal("expected RemoveSpec")
				}
			},
		},
		{
			name: "with spec replace",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				return b.WithSpecReplace(&TestSpec{})
			},
			expectedSpecs: 1,
			validate: func(t *testing.T, bp *Blueprint) {
				if _, ok := bp.Specs[0].(*ReplaceSpec); !ok {
					t.Fatal("expected ReplaceSpec")
				}
			},
		},
		{
			name: "with nested blueprint",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				nested := Blueprint{
					Name: "nested",
					Specs: []Specification{
						&TestSpec{},
						&TestSpec{},
					},
				}
				return b.WithBlueprint(nested)
			},
			expectedSpecs: 2,
			validate:      nil,
		},
		{
			name: "complex blueprint",
			buildFunc: func(b *BlueprintBuilder) *BlueprintBuilder {
				nested := Blueprint{
					Name: "nested",
					Specs: []Specification{
						&TestSpec{},
					},
				}
				return b.
					WithSpec(&TestSpec{}).
					WithSpecPresent(&TestSpec{}).
					WithBlueprint(nested).
					WithSpecRemove(&TestSpec{})
			},
			expectedSpecs: 4,
			validate:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBlueprint("test-bp")
			builder = tt.buildFunc(builder)
			bp := builder.Build()

			if bp == nil {
				t.Fatal("Build returned nil")
			}

			if len(bp.Specs) != tt.expectedSpecs {
				t.Fatalf("expected %d specs, got %d", tt.expectedSpecs, len(bp.Specs))
			}

			if tt.validate != nil {
				tt.validate(t, bp)
			}
		})
	}
}
