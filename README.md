# polyfill

[![Go report](https://goreportcard.com/badge/github.com/lofidv/polyfill)](https://goreportcard.com/report/github.com/lofidv/polyfill)

✨ **`lofidv/polyfill` is a Go library based on Go 1.18+ Generics.**

This project started as an experiment with the implementation of new generics. It may look like simple js functions but it is fully integrated and functional with arrays and structures.

As expected, generics will be much faster than implementations based on the "reflect" package. Benchmarks also show similar performance gains compared to pure `for` loops.

I feel this library is legitimate and offers many more valuable abstractions.

**Why this name?**

I wanted a **popular name**, similar to "js" and no Go package currently uses this name.

## 🚀 Install

```sh
go get github.com/lofidv/polyfill
```

This library is on beta.

## 💡 Usage

You can import `polyfill` using:

```go
import (
    "github.com/lofidv/polyfill"
)
```

Then use one of the helpers below:

```go
arrString := polyfill.NewSlice[string, int]([]string{"1", "2", "3", "4", "5"}...)

arrInt := arrString.Map(func(s string) int {
      v, _ := strconv.Atoi(s)
      return v
}) 

fmt.Println(arrInt) //[1 2 3 4 5]
```

## 🤠 Spec

GoDoc: [https://godoc.org/github.com/lofidv/polyfill](https://godoc.org/github.com/lofidv/polyfill)

Supported helpers for slices:

- Filter
- Map
- Reduce
- Some
- IndexOf
- Find

### strucs used
```go
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

```

### Filter

Iterates over a collection and returns an array of all the elements the predicate function returns `true` for.

```go
import "github.com/lofidv/polyfill"

var personsList persons
personsList = append(personsList, person{Name: "person 10", Age: 10})
personsList = append(personsList, person{Name: "person 15", Age: 15})
personsList = append(personsList, person{Name: "person 18", Age: 18})
personsList = append(personsList, person{Name: "person 20", Age: 20})


onlyAdults := polyfill.NewSlice[person, person](personsList...).Filter(func(p person) bool {
  return !p.isAdult()
})

fmt.Println(onlyAdults)// [{person 10 10} {person 15 15}]
```

### Filter with Map

Iterates over a collection and returns an array of all elements for which the predicate function returns `true`. If we add the map we can transform it to another data type

```go
s := polyfill.NewSlice[person, adolescent](personsList...)
res := s.Filter(func(p person) bool {
                return !p.isAdult()
}).Map(func(p person) adolescent {
	return adolescent{Name: p.Name, Age: p.Age}
})

fmt.Println(res) //[{person 10 10} {person 15 15}]

var p = adolescents(res)
fmt.Println(p.len()) // 2
```


### Reduce

Reduces a collection to a single value. The value is calculated by accumulating the result of running each element in the collection through an accumulator function. Each successive invocation is supplied with the return value returned by the previous call.

```go
import "github.com/lofidv/polyfill"

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

```


### Some

Returns true if at least 1 element of a subset is contained into a collection.
If the subset is empty Some returns false.

```go
sliceSomeTest := []int{1, 2, 3, 4, 5, 6}
resultSome := polyfill.NewSlice[int, int](sliceSomeTest...).Some(func(i int) bool {
	return i%2 == 0
})

fmt.Println(resultSome)// true

var sliceSomeTestPet pets
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Purin", Age: 12, Type: "dog"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Melody", Age: 1, Type: "rabbit"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Kitty", Age: 1, Type: "cat"})

resultSomePet := polyfill.NewSlice[pet, pet](sliceSomeTestPet...).Some(func(p pet) bool {
	return p.Age == 12
})

fmt.Println(resultSomePet)// true
```

### IndexOf

Returns the index at which the first occurrence of a value is found in an array or return -1 if the value cannot be found.

```go
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

fmt.Println(indexPet)//1
```

### Find

Search an element in a slice based on a predicate. It returns element and true if element was found.

```go
var sliceSomeTestPet pets
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Purin", Age: 12, Type: "dog"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Cinnamoroll", Age: 1, Type: "dog"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Melody", Age: 1, Type: "rabbit"})
sliceSomeTestPet = append(sliceSomeTestPet, pet{Name: "Kitty", Age: 1, Type: "cat"})

resultFind := polyfill.NewSlice[pet, pet](sliceSomeTestPet...).Find(func(p pet) bool {
	return p.Age == 2
})
fmt.Println(resultFind)//{ 0 }

//if results would be true p.age == 12
//{Purin 12 dog}
```


## 🤝 Contributing

- Fork the [project](https://github.com/lofidv/polyfill)
- Fix [open issues](https://github.com/lofidv/polyfill) or request new features

Don't hesitate ;)

## 👤 Authors

- jpastorm

## 💫 Show your support

Give a ⭐️ if this project helped you!

<p><a href="https://ko-fi.com/lofidev"> <img align="left" src="https://cdn.ko-fi.com/cdn/kofi3.png?v=3" height="50" width="210" alt="lofidev" /></a></p><br><br>

## 📝 License

Copyright © 2022 [jpastorm](https://github.com/lofidv).

This project is [MIT](./LICENSE) licensed.
