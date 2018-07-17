dmap is a Go library to handle dynamic objects.

It has minimal functions to get and set `map[string]interface{}`, `map[interface{}]interface{}`, and `[]interface{}`.

---

### Example

```go
package main

import (
	"fmt"

	"github.com/jeet-parekh/dmap"
)

var (
	jsonstr = `{
		"root": {
			"title": "example json",
			"contents": ["c1", "c2", "c3"]
		}
	}`
)

func main() {
	d, _ := dmap.ParseJSONBytes([]byte(jsonstr))

	// pass the keys as strings
	fmt.Println(d.Get("root", "title"))
	// output: &{example json}, <nil>

	// integers can also be passed to access elements of an array/slice
	fmt.Println(d.Get("root", "contents", 1))
	// output: &{c2}, <nil>

	// easily check if a value exists at some path
	fmt.Println(d.Exists("custom_field"))
	// output: false

	// get functions return the underlying data structures and they can be manipulated directly
	mapSI, _ := d.GetMapSI()
	mapSI["custom_field"] = "now exists"
	fmt.Println(d.Exists("custom_field"))
	// output: true

	fmt.Println(d.Exists("custom_field_2"))
	// output: false

	// values can also be set using set functions - these functions will not create the paths which do not exist - the path has to already exist
	d.SetMapSI("also exists", "custom_field_2")
	fmt.Println(d.Exists("custom_field_2"))
	// output: true

	// arrays/slices can be set the same way
	sliceI, _ := d.GetSliceI("root", "contents")
	sliceI[1] = "changed"
	fmt.Println(d.Get("root", "contents", 1))
	// output: &{changed}, <nil>
}
```
