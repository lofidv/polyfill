package main

import (
	"fmt"
	"github.com/lofidv/polyfill"
	"strconv"
)

type person struct {
	Name string
	Age  int
}

type persons []person

func (p person) isAdult() bool {
	return p.Age >= 18
}

func (p persons) len() int {
	return len(p)
}

type adolescent struct {
	Name string
	Age  int
}

type adolescents []adolescent

func (p adolescents) len() int {
	return len(p)
}

type pet struct {
	Name string
	Age  int
	Type string
}

type pets []pet

func main() {
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	s := polyfill.NewSlice[person, adolescent](personsList...)
	res := s.Filter(func(p person) bool {
		return !p.isAdult()
	}).Map(func(p person) adolescent {
		return adolescent{Name: p.Name, Age: p.Age}
	})

	fmt.Println(res) //[{person 10 10} {person 15 15}]

	var p = adolescents(res)
	fmt.Println(p.len()) // 2

	str := polyfill.NewSlice[string, int]([]string{"1", "2", "3", "4", "5"}...)
	resStr := str.Map(func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	})

	fmt.Println(resStr) //[1 2 3 4 5]

	resReduce := polyfill.NewSlice[int, int]([]int{1, 2}...).Reduce(func(acc, el int) int {
		return acc + el
	}, 0)

	fmt.Println(resReduce) //3

	var petList pets
	petList = append(petList, pet{Name: "Purin", Age: 12, Type: "dog"})
	petList = append(petList, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
	petList = append(petList, pet{Name: "Melody", Age: 1, Type: "rabbit"})
	petList = append(petList, pet{Name: "Kitty", Age: 1, Type: "cat"})

	type petMap map[string]pet

	indexed := polyfill.NewSlice[pet, petMap](petList...).Reduce(func(acc petMap, el pet) petMap {
		acc[el.Name] = el
		return acc
	}, petMap{})

	fmt.Println(indexed) ////map[Cinnamoroll:{Cinnamoroll 1 dog} Kitty:{Kitty 1 cat} Melody:{Melody 1 rabbit} Purin:{Purin 12 dog}]

	arr := []int{45, 73, 12, 98, 7, 30, 12, 85}
	index := polyfill.NewSlice[int, int](arr...).IndexOf(func(a, index int) bool {
		return a == index
	}, 12, 3)
	fmt.Println(index) //6

	var petListTwo pets
	petListTwo = append(petListTwo, pet{Name: "Purin", Age: 12, Type: "dog"})
	petListTwo = append(petListTwo, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
	petListTwo = append(petListTwo, pet{Name: "Melody", Age: 1, Type: "rabbit"})
	petListTwo = append(petListTwo, pet{Name: "Kitty", Age: 1, Type: "cat"})

	indexPet := polyfill.NewSlice[pet, string](petListTwo...).IndexOf(func(a pet, criteria string) bool {
		return a.Name == criteria
	}, "Cinnamoroll", 0)

	fmt.Println(indexPet)

	sliceSomeTest := []int{1, 2, 3, 4, 5, 6}
	resultSome := polyfill.NewSlice[int, int](sliceSomeTest...).Some(func(i int) bool {
		return i%2 == 0
	})

	fmt.Println(resultSome)

	var sliceSomeTestPet pets
	sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Purin", Age: 12, Type: "dog"})
	sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
	sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Melody", Age: 1, Type: "rabbit"})
	sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Kitty", Age: 1, Type: "cat"})

	resultSomePet := polyfill.NewSlice[pet, pet](sliceSomeTestPet...).Some(func(p pet) bool {
		return p.Age == 12
	})

	fmt.Println(resultSomePet)

	resultFind := polyfill.NewSlice[pet, pet](sliceSomeTestPet...).Find(func(p pet) bool {
		return p.Age == 2
	})
	fmt.Println(resultFind)
}
