package fiori

import (
	"fmt"
	"testing"
)

type simple struct {
	A string
	B string
	C int
	D []int
	E []string
	F []any
	G []*int
}

type complex struct {
	A string
	B string
	C int
}

func TestNew(t *testing.T) {

	f := New[simple]()

	fmt.Println(f)
}

func TestAdd(t *testing.T) {

	f := New[simple]()

	f.SetItemsPerBlock(2)

	s := simple{A: "a", B: "b", C: 3}

	f.Add(s)
	f.Add(s)
	f.Add(s)
}

func TestAddUnsafe(t *testing.T) {

	f := New[simple]()

	f.SetItemsPerBlock(2)
	f.SetUnsafeEncode(true)

	s := simple{A: "a", B: "b", C: 3}

	f.Add(s)
	f.Add(s)
	f.Add(s)
}

func TestAddComplexUnsafe(t *testing.T) {
	fmt.Println("TestAddComplexUnsafe")

	f := New[simple]()

	f.SetItemsPerBlock(2)
	f.SetUnsafeEncode(true)

	c := simple{A: "a", B: "b", C: 3, D: []int{1, 2, 3}, E: []string{"A", "B", "C"}}

	f.Add(c)
	f.Add(c)
	f.Add(c)
}
