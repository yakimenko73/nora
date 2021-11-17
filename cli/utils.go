package cli

import "github.com/mum4k/termdash/container/grid"

func addElem(e grid.Element, b *grid.Builder) {
	if e != nil {
		b.Add(e)
	}
}
