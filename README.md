# GoStream

A(nother) Stream library for Go.

Written mostly because I wanted to experiment with generics.

It uses goroutines to load objects into channels which allows for different stages of the pipeline to run in parallel.

## Streams

Streams are created with the `Slice(s []T)`, `Map(s map[K]V)`, `Of(s ...T)`, `Files(s ...string)`, `Dir(s string)` functions as well as a generic `Func(func(chan T))` function.

`Func()` takes a function that takes a channel `chan T` in which the elements are passed. It will close the stream once the function returns.

`Files()` and `Dir()` will create a Stream of interface `File` which can be satisfied by `*os.File`
A type `FileBuffer` is also provided in the package, which contains 'Reader io.ReadSeeker' as content, a `Filename string` as filename and `CloseFunc func() error` function. This can be used for files loaded into memory when, for example, stream files from a tar archive.

`Map()` will create a Stream of type `Pair[K,V]` which holds the Key and Value of the Map.

## Functions

Functions on Streams provided by the package are `Filter`, `Map` which will return a new stream with the resulting elements.
`Reduce`, `ForEach` and `Collect` will end the pipeline.

`ForEach` will just run the provided function on all elements.

`Collect` will return a Slice containing all the elements left in the pipeline.
