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
21. Close

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
2. collection.Reverse (reverse a slice)
#### Map
1. _map.Keys (get all keys of a map)
2. _map.Values (get all values of a map)
3. _map.ForEach (iterate a map)
4. _map.GetOrDefault (get value of a map by key, if the key does not exist, return the default value)
5. _map.ComputeIfAbsent (get value of a map by key, if the key does not exist, compute the value and put it into the map)
6. _map.KeysAsStream (get all keys of a map as a stream)
7. _map.ValuesAsStream (get all values of a map as a stream)
8. _map.MapAsStream (get all key-value pairs of a map as a stream)

##### Map interface
More details can be found in the [map.go](collection/map/map.go) file.
1. Put
2. PutIfAbsent
3. PutAll
4. ComputeIfAbsent
5. ComputeIfPresent
6. Get
7. GetOrDefault
8. ContainsKey
9. Keys
10. Values
11. ForEach
12. Remove
13. RemoveIf
14. Clear
15. IsEmpty
16. Size
17. AsBuiltinMap

##### HashMap
It is based on the builtin `map`, so it is not thread-safe.
```go
type HashMap[K comparable, V any] map[K]V
```
How to create a HashMap
1. `NewHashMap()`
2. `NewHashMapWithSize(size int)`
3. `NewHashMapWithMap(m Map[K]V)`
4. `NewHashMapFromBuiltinMap(m map[K]V)`

##### LinkedHashMap
It is based on the `HashMap`, so it is not thread-safe.

How to create a LinkedHashMap
1. `NewLinkedHashMap()`
2. `NewLinkedHashMapWithSize(size int)`
3. `NewLinkedHashMapWithMap(m Map[K]V)`

##### Set interface
More details can be found in the [set.go](collection/set/set.go) file.
1. Add
2. Contains
3. Remove
4. IsEmpty
5. Size
6. ForEach
7. Clear
8. AsSlice
9. Stream

##### HashSet
It is based on the builtin `map`, so it is not thread-safe.
```go
type HashSet[E comparable] map[T]struct{}
```

How to create a HashSet
1. `NewHashSet()`
2. `NewHashSetWithSize(size int)`
3. `NewHashSetWithSlice(s []E)`
4. `NewHashSetFromSet(s Set[E])`
5. `NewHashSetFromStream(s stream.Stream)`

##### LinkedHashSet
It is based on the `LinkedHashMap`, so it is not thread-safe.

How to create a LinkedHashSet
1. `NewLinkedHashSet()`
2. `NewLinkedHashSetWithSize(size int)`
3. `NewLinkedHashSetWithSlice(s []E)`
4. `NewLinkedHashSetFromSet(s Set[E])`
5. `NewLinkedHashSetFromStream(s stream.Stream)`