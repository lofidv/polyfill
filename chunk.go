package polyfill

// Chunk splits the sequence into chunks of specified size
//
// Example:
//
//	From([]int{1, 2, 3, 4, 5}).Chunk(2) // [][]int{{1, 2}, {3, 4}, {5}}
func (s *Seq[T]) Chunk(size int) [][]T {
	if size <= 0 {
		return [][]T{s.elements}
	}

	var chunks [][]T
	for i := 0; i < len(s.elements); i += size {
		end := i + size
		if end > len(s.elements) {
			end = len(s.elements)
		}
		chunks = append(chunks, s.elements[i:end])
	}
	return chunks
}
