package polyfill

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReduce(t *testing.T) {
	t.Run("testArray", testArrayReduce)
	t.Run("testStruct", testStructsReduce)
}

type testCaseReduce[V, T any] struct {
	name     string
	seed     []V
	initSeed T
	expected T
	fn       func(T, V) T
}

func testArrayReduce(t *testing.T) {
	t.Parallel()
	testsIns := []testCaseReduce[int, int]{
		{
			name:     "Convert string to int",
			seed:     []int{1, 2},
			initSeed: 0,
			expected: 3,
			fn: func(acc int, el int) int {
				return acc + el
			},
		},
	}

	for _, tc := range testsIns {
		t.Run(tc.name, runTestCaseReduce(tc))
	}

}

func testStructsReduce(t *testing.T) {
	t.Parallel()
	var petList pets
	petList = append(petList, pet{Name: "Purin", Age: 12, Type: "dog"})
	petList = append(petList, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
	petList = append(petList, pet{Name: "Melody", Age: 1, Type: "rabbit"})
	petList = append(petList, pet{Name: "Kitty", Age: 1, Type: "cat"})

	type petMap map[string]pet

	var personListExpected = petMap{
		"Cinnamoroll": pet{Name: "Cinnamoroll", Age: 1, Type: "dog"},
		"Kitty":       pet{Name: "Kitty", Age: 1, Type: "cat"},
		"Melody":      pet{Name: "Melody", Age: 1, Type: "rabbit"},
		"Purin":       pet{Name: "Purin", Age: 12, Type: "dog"},
	}

	testPersons := []testCaseReduce[pet, petMap]{
		{
			name:     "person to adolescent",
			seed:     petList,
			initSeed: petMap{},
			expected: personListExpected,
			fn: func(acc petMap, el pet) petMap {
				acc[el.Name] = el
				return acc
			},
		},
	}

	for _, tc := range testPersons {
		t.Run(tc.name, runTestCaseReduce(tc))
	}
}

func runTestCaseReduce[V, T any](tc testCaseReduce[V, T]) func(t *testing.T) {
	return func(t *testing.T) {
		res := NewSlice[V, T](tc.seed...).Reduce(tc.fn, tc.initSeed)
		assert.Equal(t, res, tc.expected)
	}
}
