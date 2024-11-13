package fiori

import (
	"fmt"
	"testing"
)

type simple struct {
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

	fmt.Println(f)

	s := simple{A: "a", B: "b", C: 3}

	f.Add(s)
	f.Add(s)
	f.Add(s)

	fmt.Println(f)
}
