# Go-DI

Go-DI implements a small and lightweight container
used for dependency injection without code generation.

Its main component is the Container interface, which can be instantiated
with NewContainer. Now dependencies can be bound to the container either
instantiating or as a singleton and queried form the Container's ResolverFunc.

It's recommended to pass the Container's ResolverFunc with the context of your
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