package _map

import "fmt"

func Example_newHashMapFromBuiltinMap() {
	lowLevelMap := map[string]int{
		"a": 1,
		"b": 2,
	}
	var m Map[string, int] = NewHashMapFromBuiltinMap[map[string]int, string, int](lowLevelMap)
	fmt.Println(m.Size())

	// Output:
	// 2
}
