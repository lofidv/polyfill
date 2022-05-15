package polyfill

type person struct {
	Name string
	Age  int
}

type persons []person

func (p person) isAdult() bool {
	return p.Age >= 18
}

type adolescent struct {
	Name string
	Age  int
}

type adolescents []adolescent

type pet struct {
	Name string
	Age  int
	Type string
}

type pets []pet
