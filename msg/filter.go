package msg

import (
	"bytes"
	"fmt"
	"sort"
)

type Filter struct {
	Included []int
	Excluded []int
}

func NewFilter() *Filter {
	return &Filter{
		Included: make([]int, 0),
		Excluded: make([]int, 0),
	}
}

func (f *Filter) In(id int) {
	f.Included = append(f.Included, id)
	sort.Ints(f.Included)
}

func (f *Filter) Out(id int) {
	f.Excluded = append(f.Excluded, id)
	sort.Ints(f.Excluded)
}

func (f *Filter) String() string {
	var buffer bytes.Buffer

	lin := len(f.Included)
	lout := len(f.Excluded)
	i := 0

	if lin > 0 {
		for _, in := range f.Included {
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
		for _, out := range f.Excluded {
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
	if f == nil {
		return true
	}

	id := m.senderId

	if len(f.Excluded) > 0 {
		for _, out := range f.Excluded {
			if out == id {
				return false
			}
		}
	}

	if len(f.Included) > 0 {
		for _, in := range f.Included {
			if in == id {
				return true
			}
		}
	} else {
		return true
	}

	return false
}
