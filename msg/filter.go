package msg

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
}

func (f *Filter) Out(id int) {
	f.out = append(f.out, id)
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
