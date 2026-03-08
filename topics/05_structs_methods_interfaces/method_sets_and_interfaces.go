package main

import "fmt"

type Stringer interface {
    String() string
}

type User struct {
    Name string
}

func (u *User) String() string {
    return "user:" + u.Name
}

func printString(s Stringer) {
    fmt.Println(s.String())
}

func main() {
    u := User{Name: "alice"}
    up := &u

    printString(up)

    var s Stringer = up
    fmt.Println(s.String())

    users := []User{{Name: "a"}, {Name: "b"}}
    for i := range users {
        printString(&users[i])
    }

    fmt.Println("A value of type User does not satisfy Stringer here because String is on *User only.")
}
