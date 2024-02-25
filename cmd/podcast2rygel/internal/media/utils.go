package media

func SliceOffsetWithMax[T any](slice []*T, offset uint, max uint) []*T {
	size := uint(len(slice))
	if offset >= size {
		return nil
	}

	end := offset + max
	if end > size {
		end = size
	}

	return slice[offset:end]
}
