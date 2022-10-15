# Go Tools
Save your life!

## Installation
```bash
go get -u github.com/carter-ya/go-tools
```

## Usage
### Stream
#### How to create a stream
##### stream.From
```go
items:= []int64{1, 2, 3, 4, 5}
s := stream.From(func(source chan<- any) {
	    for _, item := range items {
        source <- item
    }
})
```

##### stream.Just
```go
s := stream.Just([]int64{1, 2, 3, 4, 5})
```

##### stream.Range
```go
s := stream.Range[int64](0, 100)
```

##### stream.Concat
```go
s1 := stream.Just([]int64{1, 2, 3, 4, 5})
s2 := stream.Range[int64](0, 100)
s := Concat(s1, []stream.Stream{s2})
```

#### How to create a parallel stream
All the methods above can be used to create a parallel stream, 
just add `stream.WithParallelism()` to the end of the method name.
For example:
```go
s := stream.Range[int64](0, 100, stream.WithParallelism(4))
s1 := stream.Just([]int64{1, 2, 3, 4, 5}, stream.WithParallelism(4))
```

#### How to convert a parallel stream to a synchronous stream
All the methods of `stream.Stream` can be used to convert a parallel stream to a synchronous stream,
just add `stream.WithSync()` to the end of the method name.
For example:
```go
s := stream.Range[int64](0, 100, stream.WithParallelism(4)).Filter(func(item any) bool {
    return item.(int64) > 50
}, stream.WithSync())
```

#### How to use a stream
More details can be found in the [stream.go](stream/stream.go) file.
1. Map
2. FlatMap
3. Filter
4. Concat
5. Sort
6. Distinct
7. Skip
8. Limit
9. TakeWhile
10. DropWhile
11. Peek
12. AnyMatch
13. AllMatch
14. NoneMatch
15. FindFirst
16. Count
17. Reduce
18. ForEach
19. ToIfaceSlice
20. Collect
21. Done

#### How to use `Collect`
More details can be found in the [collectors.go](stream/collectors.go) file.
1. Identify
2. MapSupplier
3. MapSupplierWithSize
4. SliceSupplier
5. JoiningSupplier
6. GroupBySupplier

### Collection
#### Slice
1. collection.Shuffle (shuffle a slice)
#### Map
1. collection.Keys (get all keys of a map)
2. collection.Values (get all values of a map)
3. collection.ForEach (iterate a map)
4. collection.GetOrDefault (get value of a map by key, if the key does not exist, return the default value)
5. collection.ComputeIfAbsent (get value of a map by key, if the key does not exist, compute the value and put it into the map)
6. stream.KeysAsStream (get all keys of a map as a stream)
7. stream.ValuesAsStream (get all values of a map as a stream)
8. stream.MapAsStream (get all key-value pairs of a map as a stream)