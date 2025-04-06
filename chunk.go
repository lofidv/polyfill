package polyfill

// Chunk splits the slice into chunks of specified size
func (s *Slice[T]) Chunk(size int) []*Slice[T] {
	var chunks []*Slice[T]
	for i := 0; i < len(s.elements); i += size {
		end := i + size
		if end > len(s.elements) {
			end = len(s.elements)
		}
		chunks = append(chunks, Wrap(s.elements[i:end]))
	}
	return chunks
}
