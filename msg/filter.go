package msg

import (
	"bytes"
    "sort"
	"fmt"
)

type Filter struct {
	in  []int
	out []int
}

func NewFilter() *Filter {
	return &Filter{
		in:  make([]int, 0),
		out: make([]int, 0),
	}
}

func (f *Filter) In(id int) {
	f.in = append(f.in, id)
    sort.Ints(f.in)
}

func (f *Filter) Out(id int) {
	f.out = append(f.out, id)
    sort.Ints(f.out)
}

func (f *Filter) String() string {
	var buffer bytes.Buffer

	lin := len(f.in)
	lout := len(f.out)
	i := 0

	if lin > 0 {
		for _, in := range f.in {
			buffer.WriteString(fmt.Sprintf("%d", in))

			if i != lin-1 {
				buffer.WriteByte(',')
			}

			i++
		}

		if lout > 0 {
			buffer.WriteByte(',')
		}
	}

	i = 0

	if lout > 0 {
		for _, out := range f.out {
			buffer.WriteString(fmt.Sprintf("-%d", out))

			if i != lout-1 {
				buffer.WriteByte(',')
			}

			i++
		}
	}

	return buffer.String()
}

func (m *Message) Allowed(f *Filter) bool {
	id := m.senderId

	if len(f.out) > 0 {
		for _, out := range f.out {
			if out == id {
				return false
			}
		}
	}

	if len(f.in) > 0 {
		for _, in := range f.in {
			if in == id {
				return true
			}
		}
	} else {
		return true
	}

	return false
}
