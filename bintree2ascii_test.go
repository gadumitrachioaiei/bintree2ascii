package bintree2ascii

import (
	"encoding/json"
	"fmt"
	"testing"
)

func makeTree() *jsonNode {
	n := jsonNode{
		Name:      "1",
		LeftEdge:  "3",
		RightEdge: "2",
	}
	n.Left = &jsonNode{
		Name:      "2",
		LeftEdge:  "1",
		RightEdge: "1",
		Left:      &jsonNode{Name: "4"},
		Right:     &jsonNode{Name: "5"},
	}
	n.Right = &jsonNode{
		Name:      "3",
		RightEdge: "1",
		Right:     &jsonNode{Name: "6"},
	}
	return &n
}

func TestNode_Draw(t *testing.T) {
	at := NewAsciiArt(Config{
		NodeWidth:  4,
		NodeHeight: 1,
		EdgeHeight: 3,
		Distance:   2,
		Sep:        1,
	})
	tree := makeTree()
	s, err := json.Marshal(tree)
	if err != nil {
		t.Fatal(err)
	}
	if s, err := at.FromJson(s); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(s)
	}
	fmt.Println(at.FromInterface(jsonTree{tree}))
}
