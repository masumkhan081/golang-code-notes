package main

import (
	"context"
	"fmt"
	"log"

	"github.com/masumkhan081/golang-code-notes/topics/14_project_layout_clean_arch/internal/repository"
	"github.com/masumkhan081/golang-code-notes/topics/14_project_layout_clean_arch/internal/service"
)

func main() {
	// 1. Initialize dependencies (Dependency Injection)
	repo := repository.NewInMemoryUserRepo()
	svc := service.NewUserService(repo)

	ctx := context.Background()

	// 2. Use the service
	fmt.Println("Registering user...")
	err := svc.RegisterUser(ctx, "1", "Alice", "alice@example.com", 30)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Retrieve user
	user, err := svc.GetUser(ctx, "1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User found: %+v\n", user)
}
