package messages

func valOrZero[T any](p *T) T {
	var v T
	if p != nil {
		v = *p
	}
	return v
}
