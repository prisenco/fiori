package fiori

import (
	"fmt"
	"testing"
)

type simple struct {
	a string
	b string
	c int
}

func TestNew(t *testing.T) {

	f := New[simple]()

	fmt.Println(f)
}

func TestAdd(t *testing.T) {

	f := New[simple]()

	f.SetItemsPerBlock(2)

	fmt.Println(f)

	s := simple{a: "a", b: "b", c: 3}

	f.Add(s)
	f.Add(s)
	f.Add(s)

	fmt.Println(f)
}
