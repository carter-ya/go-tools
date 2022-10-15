package stream

import "fmt"

func ExampleStream_Map() {
	p := []struct {
		Name string
		Age  int
	}{
		{"Alice", 20},
		{"Bob", 30},
	}
	Just(p).
		Map(func(item any) any {
			return item.(struct {
				Name string
				Age  int
			}).Name
		}).
		ForEach(func(item any) {
			fmt.Println("name:", item.(string))
		})
}
