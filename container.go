// Package godi implements a small and lightweight container
// used for dependency injection without code generation.
//
// Its main component is the Container interface, which can be instantiated
// with NewContainer. Now dependencies can be bound to the container either
// instantiating or as a singleton and queried form the Container's ResolverFunc.
//
// It's recommended to pass the Container's ResolverFunc with the context of your
// application, allowing you to access the dependency injection container easily.
//
//	// create the container and bind a dependency
//	container := godi.NewContainer()
//	container.MustBind("time-service", func(resolver godi.ResolverFunc) any {
//		return time.Now()
//	})
//
//	// add resolver to a context
//	ctx := context.WithValue(context.Background(), "container", container.Resolver())
//
//	// use the context to retrieve the resolver and execute it
//	resolver := ctx.Value("container").(godi.ResolverFunc)
//	currentTime := godi.MustResolve[time.Time]("time-service", resolver)
//	fmt.Println(currentTime.Unix())
package godi

import (
	"errors"
	"fmt"
	"sync"
)

// ResolverFunc is a generic function, used to request a dependency from
// a Container by its name. As the returned value is of any value, you may
// use the Resolve or MustResolve helper functions to handle the type
// conversion for you.
type ResolverFunc = func(string) (any, error)

// BinderFunc is a generic function, used to bind dependencies to a Container.
// It's first argument is a ResolverFunc, which allows you to request additional
// dependencies as needed.
type BinderFunc = func(resolver ResolverFunc) any

// Container is the main interface for the dependency collection container.
// Through the Container, multiple dependencies can be prepared and stored
// by an identifying name and resolved on demand by this name.
//
// The Container supports instanced binding, through its Bind method.
// Instanced dependencies are instanced on demand, if the dependency is
// requested. The Container also supports singleton binding, through its
// BindSingleton method. Singleton dependencies are instanced once lazily,
// when requested for the first time. All further dependency requests
// receive this first instance. Both binding methods offer a variant, which
// panics on a failed bind.
//
// Once all Dependencies are bound to the container. You may call Lock
// to prevent any more modification of the allowed dependencies. To resolve
// a dependency by its name, get the ResolverFunc by calling Resolver. You
// may use the Resolve or MustResolve helper functions to handle the type
// conversion for you.
type Container interface {
	Lock()
	Bind(name string, binder BinderFunc) error
	MustBind(name string, binder BinderFunc)
	BindSingleton(name string, binder BinderFunc) error
	MustBindSingleton(name string, binder BinderFunc)
	Resolver() ResolverFunc
}

// NewContainer instantiates a generic Container, which can be filled
// with instanced or singleton dependencies, locked and queried for
// dependencies.
func NewContainer() Container {
	s := defaultContainer{
		locked:   false,
		services: make(map[string]BinderFunc),
	}
	return &s
}

type defaultContainer struct {
	locked   bool
	services map[string]BinderFunc
}

func (d *defaultContainer) Lock() {
	d.locked = true
}

func (d *defaultContainer) Bind(name string, binder BinderFunc) error {
	if d.locked {
		return errors.New("service container locked. no more services can be bound")
	}
	if _, ok := d.services[name]; ok {
		return errors.New(fmt.Sprintf("service with name %s already bound", name))
	}
	d.services[name] = binder
	return nil
}

func (d *defaultContainer) MustBind(name string, binder BinderFunc) {
	if err := d.Bind(name, binder); err != nil {
		panic(err.Error())
	}
}

func (d *defaultContainer) BindSingleton(name string, binder BinderFunc) error {
	var lazyBind sync.Once
	var result any
	bind := func(resolver ResolverFunc) any {
		lazyBind.Do(func() {
			result = binder(resolver)
		})
		return result
	}
	return d.Bind(name, bind)
}

func (d *defaultContainer) MustBindSingleton(name string, binder BinderFunc) {
	if err := d.BindSingleton(name, binder); err != nil {
		panic(err.Error())
	}
}

func (d *defaultContainer) Resolver() ResolverFunc {
	return func(name string) (any, error) {
		if _, ok := d.services[name]; !ok {
			return nil, errors.New(fmt.Sprintf("%s service not found in container", name))
		}
		return d.services[name](d.Resolver()), nil
	}
}
