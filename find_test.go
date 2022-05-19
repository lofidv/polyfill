package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	t.Run("testArray", testArrayFind)
	t.Run("testStruct", testStructsFind)
}

type testCaseFind[V, T any] struct {
	name     string
	seed     []V
	expected V
	fn       func(V) bool
}

func testArrayFind(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseFind[int, int]{
		{
			name:     "i == 1",
			seed:     []int{1, 2, 3, 4},
			expected: 1,
			fn:       func(i int) bool { return i == 1 },
		},
		{
			name:     "i == 6",
			seed:     []int{1, 2, 3, 4},
			expected: 0,
			fn:       func(i int) bool { return i == 6 },
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseFind(tc))
	}

	testsStrings := []testCaseFind[string, int]{
		{
			name:     "i == a",
			seed:     []string{"a", "b", "c", "d"},
			expected: "a",
			fn:       func(a string) bool { return a == "a" },
		},
	}

	for _, tc := range testsStrings {
		t.Run(tc.name, runTestCaseFind(tc))
	}

}

func testStructsFind(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	testPersons := []testCaseFind[person, person]{
		{
			name: "p.age == 18",
			seed: personsList,
			expected: person{
				Name: "person 18",
				Age:  18,
			},
			fn: func(p person) bool { return p.Age == 18 },
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseFind(tc))
	}
}

func runTestCaseFind[V, T any](tc testCaseFind[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).Find(tc.fn)
		assert.Equal(t, res, tc.expected)
	}
}
