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

type Adolescent struct {
	Name string
	Age  int
}

type Adolescents []Adolescent

func (p Adolescents) len() int {
	return len(p)
}

func main() {
	var personsList persons
	personsList = append(personsList, person{Name: "person 10", Age: 10})
	personsList = append(personsList, person{Name: "person 15", Age: 15})
	personsList = append(personsList, person{Name: "person 18", Age: 18})
	personsList = append(personsList, person{Name: "person 20", Age: 20})

	s := polyfill.NewSlice[person, Adolescent](personsList...)
	res := s.Filter(func(p person) bool {
		return !p.isAdult()
	}).Map(func(p person) Adolescent {
		return Adolescent{Name: p.Name, Age: p.Age}
	})

	fmt.Println(res)
	var p = Adolescents(res)
	fmt.Println(p.len())

	str := polyfill.NewSlice[string, int]([]string{"1", "2", "3", "4", "5"}...)

	resStr := str.Map(func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	})

	fmt.Println(resStr) //[1 2 3 4 5]
}
