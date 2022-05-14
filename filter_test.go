package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("testArray", testArrayFilter)
	t.Run("testStruct", testStructsFilter)
}

type testCaseFilter[V, T comparable] struct {
	name     string
	seed     []V
	expected Slice[V, T]
	fn       func(V) bool
}

func testArrayFilter(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseFilter[int, int]{
		{
			name:     "i == 1",
			seed:     []int{1, 2, 3, 4},
			expected: []int{1},
			fn:       func(i int) bool { return i == 1 },
		},
		{
			name:     "i % 2 == 0",
			seed:     []int{1, 2, 3, 4},
			expected: []int{2, 4},
			fn:       func(i int) bool { return i%2 == 0 },
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseFilter(tc))
	}

	testsStrings := []testCaseFilter[string, int]{
		{
			name:     "i != a",
			seed:     []string{"a", "b", "c", "d"},
			expected: []string{"b", "c", "d"},
			fn:       func(a string) bool { return a != "a" },
		},
	}

	for _, tc := range testsStrings {
		t.Run(tc.name, runTestCaseFilter(tc))
	}

}

func testStructsFilter(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	var personListExpected = Slice[person, person]{
		person{Name: "person 18", Age: 18},
		person{Name: "person 20", Age: 20},
	}

	testPersons := []testCaseFilter[person, person]{
		{
			name:     "p.age >= 18",
			seed:     personsList,
			expected: personListExpected,
			fn:       func(p person) bool { return p.isAdult() },
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseFilter(tc))
	}
}

func runTestCaseFilter[V, T comparable](tc testCaseFilter[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).Filter(tc.fn)
		assert.Equal(t, res, tc.expected)
	}
}
