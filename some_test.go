package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSome(t *testing.T) {
	t.Run("testArray", testArraySome)
	t.Run("testStruct", testStructsSome)
}

type testCaseSome[V, T any] struct {
	name     string
	seed     []V
	expected bool
	fn       func(V) bool
}

func testArraySome(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseSome[int, int]{
		{
			name:     "are there values even?",
			seed:     []int{1, 2, 3, 4, 5, 6},
			expected: true,
			fn:       func(a int) bool { return a%2 == 0 },
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseSome(tc))
	}
}

func testStructsSome(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	testPersons := []testCaseSome[person, person]{
		{
			name:     "Is there someone with 25 years?",
			seed:     personsList,
			expected: false,
			fn:       func(a person) bool { return a.Age == 25 },
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseSome(tc))
	}
}

func runTestCaseSome[V, T any](tc testCaseSome[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).Some(tc.fn)
		assert.Equal(t, res, tc.expected)
	}
}
