package main

type Productor struct {
	items   [][]string
	indexes []int
	ch      chan []int
}

func NewProductor(items [][]string, ch chan []int) *Productor {
	return &Productor{
		items:   items,
		indexes: make([]int, len(items)),
		ch:      ch,
	}
}

func (p *Productor) FindNext(index_i int) {
	if index_i == len(p.items) {
		indexes := make([]int, len(p.indexes))
		copy(indexes, p.indexes)
		p.ch <- indexes
		return
	}

	for i := 0; i < len(p.items[index_i]); i++ {
		p.indexes[index_i] = i
		p.FindNext(index_i + 1)
	}
}

func Product(items [][]string) chan []int {
	ch := make(chan []int, 16)
	go func() {
		p := NewProductor(items, ch)
		p.FindNext(0)
		close(p.ch)
	}()
	return ch
}
