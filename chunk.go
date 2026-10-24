package polyfill

// Chunk splits the sequence into consecutive chunks of size elements.
// The last chunk may be shorter.
//
// Chunk panics when size is less than 1: a zero size can never terminate and a
// negative size is a programming error, so failing loudly beats an infinite
// loop.
//
// The returned chunks share the receiver's backing array. Treat them as
// read-only, or Clone a chunk before modifying it.
//
//	From([]int{1, 2, 3, 4, 5}).Chunk(2) // [[1 2] [3 4] [5]]
func (s *Seq[T]) Chunk(size int) [][]T {
	if size < 1 {
		panic("polyfill: Chunk size must be greater than zero")
	}
	if s == nil || s.err != nil || len(s.elements) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(s.elements)+size-1)/size)
	for i := 0; i < len(s.elements); i += size {
		end := i + size
		if end > len(s.elements) {
			end = len(s.elements)
		}
		chunks = append(chunks, s.elements[i:end])
	}
	return chunks
}
