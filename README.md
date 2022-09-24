# Go-DI

Go-DI implements a small and lightweight container
used for dependency injection without code generation.

Its main component is the `Container` interface, which can be instantiated
with `godi.NewContainer()`. Now dependencies can be bound to the container either
instantiating or as a singleton and queried from the Container's `ResolverFunc`.

It's recommended to pass the Container's `ResolverFunc` with the context of your
application, allowing you to access the dependency injection container easily.

```go
// create the container and bind a dependency
container := godi.NewContainer()
container.MustBind("time-service", func(resolver godi.ResolverFunc) any {
    return time.Now()
})

// add resolver to a context
ctx := context.WithValue(context.Background(), "container", container.Resolver())

// use the context to retrieve the resolver and execute it
resolver := ctx.Value("container").(godi.ResolverFunc)
currentTime := godi.MustResolve[time.Time]("time-service", resolver)

fmt.Println(currentTime.Unix())
```

## Types of bound dependencies
Go-DI differentiates two types of dependencies: instantiating and singleton
dependencies.

### Instantiating Dependencies
Instantiating Dependencies are instantiated every time, they are requested
from the dependency container. Every request form the dependency container
will therefore yield a new instance.

````go
container.MustBind("rng", func(resolver godi.ResolverFunc) any {
    return rand.Int63()
})
resolver := container.Resolver()

// Will yield different results for all DI requests
fmt.Println(godi.MustResolve[int64]("rng", resolver))
fmt.Println(godi.MustResolve[int64]("rng", resolver))
fmt.Println(godi.MustResolve[int64]("rng", resolver))
````

### Singleton Dependencies
Singleton Dependencies are instantiated only once. All subsequent requests
to the dependency container will yield the first and only instantiated instance.

````go
container.MustBindSingleton("rng-once", func(resolver godi.ResolverFunc) any {
    return rand.Int63()
})
resolver := container.Resolver()

// Will yield the same result for all DI requests
fmt.Println(godi.MustResolve[int64]("rng-once", resolver))
fmt.Println(godi.MustResolve[int64]("rng-once", resolver))
fmt.Println(godi.MustResolve[int64]("rng-once", resolver))
````
