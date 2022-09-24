package godi

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultContainer_Bind(t *testing.T) {
	container := NewContainer()
	handler := func(resolver ResolverFunc) any {
		return 12345
	}
	list := []string{"foo", "bar", "baz"}
	for _, name := range list {
		err := container.Bind(name, handler)
		if err != nil {
			t.Fatalf("Unable to instanced bind dependency %s to default container", name)
		}
	}
	err := container.Bind("foo", handler)
	if err == nil {
		t.Fatalf("Could override already existing dependency %s", "foo")
	}
}

func TestDefaultContainer_MustBind(t *testing.T) {
	container := NewContainer()
	handler := func(resolver ResolverFunc) any {
		return true
	}
	container.MustBind("foo", handler)
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBind did not panic, when it should have")
		}
	}()
	container.MustBind("foo", handler)
}

func TestDefaultContainer_Bind2(t *testing.T) {
	container := NewContainer()
	container.MustBind("time-service", func(resolver ResolverFunc) any {
		return time.Now()
	})
	currentTime := MustResolve[time.Time]("time-service", container.Resolver())
	fmt.Println(currentTime.Unix())
}

func TestDefaultContainer_BindSingleton(t *testing.T) {
	container := NewContainer()
	handler := func(resolver ResolverFunc) any {
		return 12345
	}
	list := []string{"foo", "bar", "baz"}
	for _, name := range list {
		err := container.BindSingleton(name, handler)
		if err != nil {
			t.Fatalf("Unable to instanced bind dependency %s to default container", name)
		}
	}
	err := container.BindSingleton("foo", handler)
	if err == nil {
		t.Fatalf("Could override already existing dependency %s", "foo")
	}
}

func TestDefaultContainer_MustBindSingleton(t *testing.T) {
	container := NewContainer()
	handler := func(resolver ResolverFunc) any {
		return true
	}
	container.MustBindSingleton("foo", handler)
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBind did not panic, when it should have")
		}
	}()
	container.MustBindSingleton("foo", handler)
}

func TestDefaultContainer_Resolver(t *testing.T) {
	container := NewContainer()
	container.MustBind("counter", func(resolver ResolverFunc) any {
		var num uint8 = 9
		return num
	})

	resolver := container.Resolver()
	rawResult, err := resolver("counter")
	if err != nil {
		t.Fatalf("Could not resolve existing dependency %s", "counter")
	}
	result, ok := rawResult.(uint8)
	if !ok {
		t.Fatalf("Resolved dependency not of the expected type")
	}
	if result != 9 {
		t.Fatalf("Resolved dependency not the expected value. Got %d expected %d", result, 9)
	}

	_, err = resolver("foobar")
	if err == nil {
		t.Fatalf("Resolved dependency for non existing name %s", "foobar")
	}
}

func TestDefaultContainer_Resolver_Instanced(t *testing.T) {
	container := NewContainer()
	var num = 10
	container.MustBind("rand", func(resolver ResolverFunc) any {
		num *= 10
		value := num
		return value
	})
	a := MustResolve[int]("rand", container.Resolver())
	b := MustResolve[int]("rand", container.Resolver())
	if a == b {
		t.Fatalf("Expected different results, got same results. %d, %d", a, b)
	}
}

func TestDefaultContainer_Resolver_Singleton(t *testing.T) {
	container := NewContainer()
	var num = 10
	container.MustBindSingleton("rand", func(resolver ResolverFunc) any {
		num *= 10
		value := num
		return value
	})
	a := MustResolve[int]("rand", container.Resolver())
	b := MustResolve[int]("rand", container.Resolver())
	if a != b {
		t.Fatalf("Expected the same result, got different results. %d, %d", a, b)
	}
}

func TestDefaultContainer_Lock(t *testing.T) {
	handler := func(resolver ResolverFunc) any {
		return true
	}
	container := NewContainer()
	container.MustBind("foo", handler)
	container.Lock()
	err := container.Bind("bar", handler)
	if err == nil {
		t.Fatalf("Dependency can be pushed to locked container")
	}
}
