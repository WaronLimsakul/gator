package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/WaronLimsakul/gator/internal/config"
	"github.com/WaronLimsakul/gator/internal/database"

	// need to do this because we will not use it directly
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	// basic the Queries type has all method which are SQL function
	db *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

// will register the command by linked its name to the function via the command map
func (c *commands) registerCommand(name string, f func(*state, command) error) {
	if len(name) == 0 || f == nil {
		return
	}
	c.cmdMap[name] = f
	return
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmdMap[cmd.name]
	if !ok {
		return fmt.Errorf("command '%s' not found", cmd.name)
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

	gatorCommands.registerCommand("login", handlerLogin)
	gatorCommands.registerCommand("register", handlerRegister)
	gatorCommands.registerCommand("reset", handlerReset)
	gatorCommands.registerCommand("users", handlerUsersList)
	gatorCommands.registerCommand("agg", handlerAggregator)
	gatorCommands.registerCommand("addfeed", middlewareLoggedIn(handlerAddFeed))
	gatorCommands.registerCommand("feeds", handlerFeedsList)
	gatorCommands.registerCommand("follow", middlewareLoggedIn(handlerFollow))
	gatorCommands.registerCommand("following", middlewareLoggedIn(handlerFollowsList))
	gatorCommands.registerCommand("unfollow", middlewareLoggedIn(handlerUnfollow))

	// open database connection
	db, err := sql.Open("postgres", curState.config.DbUrl)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// create a queries instance
	dbQueries := database.New(db)

	// set those query to be accessible within the current state
	curState.db = dbQueries

	// os.Args first element is the location where these are store (/tmp/something)
	if len(os.Args) < 2 {
		fmt.Println("command required")
		os.Exit(1)
	}

	// first is name, and the rest are args
	inputCmd := command{os.Args[1], os.Args[2:]}

	err = gatorCommands.run(&curState, inputCmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
