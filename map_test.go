package polyfill

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("testArray", testArrayMap)
	t.Run("testStruct", testStructsMap)
}

type testCaseMap[V, T any] struct {
	name     string
	seed     []V
	expected Slice[T, V]
	fn       func(V) T
}

func testArrayMap(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseMap[string, int]{
		{
			name:     "Convert string to int",
			seed:     []string{"1", "2", "3", "4", "5"},
			expected: []int{1, 2, 3, 4, 5},
			fn: func(i string) int {
				v, _ := strconv.Atoi(i)
				return v
			},
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseMap(tc))
	}

}

func testStructsMap(t *testing.T) {
	t.Parallel()
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	var personListExpected = Slice[adolescent, person]{
		adolescent{Name: "person 10", Age: 10},
		adolescent{Name: "person 15", Age: 15},
		adolescent{Name: "person 18", Age: 18},
		adolescent{Name: "person 20", Age: 20},
	}

	testPersons := []testCaseMap[person, adolescent]{
		{
			name:     "person to adolescent",
			seed:     personsList,
			expected: personListExpected,
			fn: func(i person) adolescent {
				return adolescent{Name: i.Name, Age: i.Age}
			},
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseMap(tc))
	}
}

func runTestCaseMap[T, V any](tc testCaseMap[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).Map(tc.fn)
		assert.Equal(t, res, tc.expected)
	}
}
