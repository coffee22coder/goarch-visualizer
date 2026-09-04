package testdata

import "fmt"

type User struct {
	name string
}

type Model struct {
	id   int
	name string
}

func test() {
	u := User{name: "Jenya"}
	m := Model{id: 23, name: "QWERTY"}

	SayHello(u)
	PrintModel(m)

}

func SayHello(user User) string {
	return fmt.Sprintf("Hello, %s", user.name)
}

func PrintModel(model Model) {
	fmt.Println(model)
}
