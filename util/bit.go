package util

func BitOn(b int, index int) bool {
	return (1<<index)&b != 0
}

func ClearBit(b *int, index int) {
	*b &= ^(1 << index)
}

func SetBit(b *int, index int) {
	*b |= (1 << index)
}
