package main

type productor struct {
	items   [][]string
	indexes []int
	ch      chan []int
}

func newProductor(items [][]string, ch chan []int) *productor {
	return &productor{
		items:   items,
		indexes: make([]int, len(items)),
		ch:      ch,
	}
}

func (p *productor) findNext(index_i int) {
	if index_i == len(p.items) {
		indexes := make([]int, len(p.indexes))
		copy(indexes, p.indexes)
		p.ch <- indexes
		return
	}

	for i := 0; i < len(p.items[index_i]); i++ {
		p.indexes[index_i] = i
		p.findNext(index_i + 1)
	}
}

func Product(items [][]string) chan []int {
	ch := make(chan []int, 16)
	go func() {
		p := newProductor(items, ch)
		p.findNext(0)
		close(p.ch)
	}()
	return ch
}
