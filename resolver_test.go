package godi

import (
	"testing"
)

func TestResolve(t *testing.T) {
	var err error
	container := NewContainer()
	container.MustBind("foo", func(resolver ResolverFunc) any {
		return 1
	})
	container.MustBind("bar", func(resolver ResolverFunc) any {
		return 2
	})

	fooV, err := Resolve[int]("foo", container.Resolver())
	if err != nil {
		t.Fatalf("Dependency %s not bound", "foo")
	}
	if fooV != 1 {
		t.Fatalf("Dependency %s has unexpected value. Expected %d got %d", "foo", 1, fooV)
	}

	barVstr, err := Resolve[string]("bar", container.Resolver())
	if err == nil {
		t.Fatalf("Dependency resolved with wrong type. Expected error but got string value %s", barVstr)
	}

	barV, err := Resolve[int]("bar", container.Resolver())
	if err != nil {
		t.Fatalf("Dependency %s not bound", "bar")
	}
	if barV != 2 {
		t.Fatalf("Dependency %s has unexpected value. Expected %d got %d", "bar", 2, barV)
	}

	_, err = Resolve[int]("baz", container.Resolver())
	if err == nil {
		t.Fatalf("Unexpected resolving of non existing dependency %s", "baz")
	}
}

func TestMustResolve(t *testing.T) {
	container := NewContainer()
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustResolve did not panic, when it should have")
		}
	}()
	MustResolve[int]("test", container.Resolver())
}
