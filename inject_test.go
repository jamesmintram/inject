package inject_test

import (
	"testing"

	"github.com/jamesmintram/inject"
)

type StructB struct {
}

type StructA struct {
	AB *StructB `inject:""`
}

type StructC struct {
	AB *StructB `inject:""`
}

func TestSimple(t *testing.T) {
	var ToInjectA struct {
		B *StructB `inject:""`
	}

	inject.Inject(&ToInjectA)

	if ToInjectA.B == nil {
		t.Error("Injected with different instances")
	}
}

func TestMultiple(t *testing.T) {
	var ToInjectA struct {
		A *StructA `inject:""`
		B *StructB `inject:""`
	}

	inject.Inject(&ToInjectA)

	if ToInjectA.A == nil {
		t.Error("Injected with different instances")
	}
	if ToInjectA.B == nil {
		t.Error("Injected with different instances")
	}
}

func TestTransitive(t *testing.T) {
	var ToInjectA struct {
		A *StructA `inject:""`
	}

	inject.Inject(&ToInjectA)

	if ToInjectA.A.AB == nil {
		t.Error("Injected with different instances")
	}
}

func TestMultipleTransitive(t *testing.T) {
	var ToInjectA struct {
		A *StructA `inject:""`
		C *StructC `inject:""`
	}

	inject.Inject(&ToInjectA)

	if ToInjectA.A == nil {
		t.Error("Injected with different instances")
	}
	if ToInjectA.C == nil {
		t.Error("Injected with different instances")
	}
}

func TestMultipleTransitiveSameInstances(t *testing.T) {
	var ToInjectA struct {
		A *StructA `inject:""`
		C *StructC `inject:""`
	}

	inject.Inject(&ToInjectA)

	if ToInjectA.A.AB != ToInjectA.C.AB {
		t.Error("Injected with different instances")
	}
}

func TestInjectAllSameInstances(t *testing.T) {
	var ToInjectA struct {
		A *StructA `inject:""`
		B *StructB `inject:""`
	}
	var ToInjectB struct {
		A *StructA `inject:""`
		B *StructB `inject:""`
	}

	inject.InjectAll(&ToInjectA, &ToInjectB)

	if ToInjectA.A != ToInjectB.A {
		t.Error("Injected with different instances")
	}
	if ToInjectA.B != ToInjectB.B {
		t.Error("Injected with different instances")
	}
}

func TestInjectUseExisting(t *testing.T) {
	var testA = StructA{}
	var ToInjectB struct {
		A *StructA `inject:""`
		B *StructB `inject:""`
	}

	inject.InjectAll(&testA, &ToInjectB)

	if ToInjectB.A != &testA {
		t.Error("Not using services supplied in a map")
	}
}
