package godi

import (
	"errors"
	"fmt"
)

// MustResolve is a helper function to simplify interaction with a
// ResolverFunc. MustResolve tries to fetch a dependency by its name
// and panics, if the dependency can't be converted to the given type
// or can't be found by the provided ResolverFunc.
func MustResolve[T any](name string, resolver ResolverFunc) T {
	value, err := Resolve[T](name, resolver)
	if err != nil {
		panic(err)
	}
	return value
}

// Resolve is a helper function to simplify interaction with a ResolverFunc.
// Resolve tries to fetch a dependency by its name and convert it to the given
// type. An error is returned if the conversion failed or the dependency could
// not be found.
func Resolve[T any](name string, resolver ResolverFunc) (T, error) {
	t, err := resolver(name)
	if err != nil {
		var res T
		return res, err
	}
	v, ok := t.(T)
	if !ok {
		return v, errors.New(fmt.Sprintf("Unable to convert %s to the requested type", name))
	}
	return v, nil
}
