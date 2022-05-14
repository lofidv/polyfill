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
