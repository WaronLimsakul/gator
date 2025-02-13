package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WaronLimsakul/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	loginName := cmd.args[0]

	if _, err := s.db.GetUser(context.Background(), loginName); err != nil {
		fmt.Println("username not found")
		os.Exit(1)
	}

	if err := s.config.SetUser(loginName); err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		fmt.Println("name not found")
		return fmt.Errorf("name not found")
	}

	requestedName := cmd.args[0]

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      requestedName,
	}

	// only thing that can go wrong here is that user name is taken
	createdUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		fmt.Printf("user name '%s' is already taken\n", requestedName)
		os.Exit(1)
	}

	fmt.Printf("User '%s' registered\n", requestedName)

	if err := s.config.SetUser(requestedName); err != nil {
		return err
	}

	fmt.Println("User set. user's data:")
	fmt.Printf("%v", createdUser)

	return nil
}

func handlerReset(s *state, cmd command) error {
	// only have user now
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database reset")
	return nil
}

func handlerUsersList(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
