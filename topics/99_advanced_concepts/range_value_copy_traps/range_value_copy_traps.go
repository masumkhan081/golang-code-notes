package main

import "fmt"

type User struct {
    Name  string
    Score int
}

func wrongMutation(users []User) {
    for _, u := range users {
        u.Score += 10
    }
}

func correctMutationByIndex(users []User) {
    for i := range users {
        users[i].Score += 10
    }
}

func pointerTrap() []*int {
    var out []*int
    for i := 0; i < 3; i++ {
        v := i
        out = append(out, &v)
    }
    return out
}

func main() {
    users := []User{
        {Name: "alice", Score: 10},
        {Name: "bob", Score: 20},
    }

    wrongMutation(users)
    fmt.Println("after wrongMutation:", users)

    correctMutationByIndex(users)
    fmt.Println("after correctMutationByIndex:", users)

    ptrs := pointerTrap()
    for _, p := range ptrs {
        fmt.Println(*p)
    }

    m := map[string]int{"a": 1, "b": 2, "c": 3}
    for k, v := range m {
        fmt.Println("map iteration:", k, v)
    }
}
