package main

import (
	"fmt"
	"os"

	"github.com/WaronLimsakul/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	loginName := cmd.args[0]
	if err := s.config.SetUser(loginName); err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if len(name) == 0 || f == nil {
		return
	}
	c.cmdMap[name] = f
	return
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmdMap[cmd.name]
	if !ok {
		return fmt.Errorf("'%s' not found", cmd.name)
	}

	if err := handler(s, cmd); err != nil {
		return err
	}
	return nil
}

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		return
	}
	curState := state{config: &conf}

	gatorCommands := commands{}
	gatorCommands.cmdMap = make(map[string]func(*state, command) error)
	gatorCommands.register("login", handlerLogin)

	// os.Args first element is the location where these are store (/tmp/something)
	if len(os.Args) < 2 {
		fmt.Println("command required")
		os.Exit(1)
	}
	inputCmd := command{os.Args[1], os.Args[2:]}
	err = gatorCommands.run(&curState, inputCmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
