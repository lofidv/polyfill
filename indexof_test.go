package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndex(t *testing.T) {
	t.Run("testArray", testArrayIndexOf)
	t.Run("testStruct", testStructsIndexOf)
}

type testCaseIndexOf[V, T any] struct {
	name     string
	seed     []V
	expected int
	fn       func(V, T) bool
	value    T
	pos      int
}

func testArrayIndexOf(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseIndexOf[int, int]{
		{
			name:     "search the value 12 of array starting at 3",
			seed:     []int{45, 73, 12, 98, 7, 30, 12, 85},
			expected: 6,
			fn:       func(a, index int) bool { return a == index },
			value:    12,
			pos:      3,
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseIndexOf(tc))
	}
}

func testStructsIndexOf(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	testPersons := []testCaseIndexOf[person, string]{
		{
			name:     "search the person 18 of the array person starting in 0",
			seed:     personsList,
			expected: 2,
			fn:       func(a person, criteria string) bool { return a.Name == criteria },
			value:    "person 18",
			pos:      0,
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseIndexOf(tc))
	}
}

func runTestCaseIndexOf[V, T any](tc testCaseIndexOf[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).IndexOf(tc.fn, tc.value, tc.pos)
		assert.Equal(t, res, tc.expected)
	}
}
