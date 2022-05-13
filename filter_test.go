package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {
	t.Run("testArray", testArray)
	t.Run("testStruct", testStructs)
}

type testCase[V Constraint] struct {
	name     string
	seed     []V
	expected []V
	fn       func(V) bool
}

func testArray(t *testing.T) {
	t.Parallel()
	testsIns := []testCase[int]{
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
		t.Run(tc.name, runTestCase(tc))
	}

	testsStrings := []testCase[string]{
		{
			name:     "i != a",
			seed:     []string{"a", "b", "c", "d"},
			expected: []string{"b", "c", "d"},
			fn:       func(a string) bool { return a != "a" },
		},
	}

	for _, tc := range testsStrings {
		t.Run(tc.name, runTestCase(tc))
	}

}

type person struct {
	Name string
	Age  int
}

type persons []person

func (p person) isAdult() bool {
	return p.Age >= 18
}

func testStructs(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	var personListExpected = persons{
		person{
			Name: "person 18",
			Age:  18,
		},
		person{
			Name: "person 20",
			Age:  20,
		},
	}

	testPersons := []testCase[person]{
		{
			name:     "p.age >= 18",
			seed:     personsList,
			expected: personListExpected,
			fn:       func(p person) bool { return p.isAdult() },
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCase(tc))
	}
}

func runTestCase[V Constraint](tc testCase[V]) func(t *testing.T) {
	return func(t *testing.T) {
		var s Slice[V]
		res := s.Add(tc.seed...).Filter(tc.fn)

		assert.Equal(t, res, tc.expected)
	}
}
