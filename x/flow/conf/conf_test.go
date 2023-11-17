package conf

import (
	"testing"
	"fmt"


	"gopkg.in/yaml.v3"
)

//--------------------------------------------------------------------------------

type MyData struct {
	A int
	B string
	C []any
	D any
}

type MyData2 struct {
	Name string
}

type MyData3 struct {
	Number int
}

func TestThing(t *testing.T) {
	registry.Clear()
	Register("MyData", MyData{})
	Register("MyData2", MyData2{})
	Register("MyData3", MyData3{})

	data := `
A: 5
B: "hello"
C:
  - {MyData2: {Name: "slicestring"}}
  - {MyData3: {Number: 77}}
D: {MyData2: {Name: "secondstring"}}
`

	m := make(map[string]any)
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil { panic(err) }
	printMap(m)

	result := MyData{}
	err = Decode(m, &result)
	if err != nil { panic(err) }

	fmt.Println("Result:", result)
	fmt.Printf("Type: %T\n", result.D)
}

func TestMerge(t *testing.T) {
	registry.Clear()
	Register("MyData", MyData{})
	Register("MyData2", MyData2{})
	Register("MyData3", MyData3{})

	data := `
A: 5
B: "hello"
C:
  - {MyData2: {Name: "slicestring"}}
  - {MyData3: {Number: 77}}
D: {MyData2: {Name: "secondstring"}}
`

	data2 := `
B: "hellothere-overwrittenbydata2"
D: {MyData2: {Name: "overwrittenbydata2"}}
`

	m := make(map[string]any)
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil { panic(err) }
	fmt.Println("Before")
	printMap(m)

	yaml.Unmarshal([]byte(data2), &m)
	fmt.Println("After")
	printMap(m)

	result := MyData{}
	err = Decode(m, &result)
	if err != nil { panic(err) }

	fmt.Println("Result:", result)
	fmt.Printf("Type: %T\n", result.D)
}

func printMap(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			fmt.Printf(k + " . ")
			printMap(t)
			fmt.Printf("\n")
		default:
			fmt.Println(k, ":", v)
		}
	}
}

