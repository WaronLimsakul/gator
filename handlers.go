package main

import (
	"context"
	"fmt"
	"os"

	"github.com/WaronLimsakul/gator/internal/database"
	"github.com/WaronLimsakul/gator/internal/rss"
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
		ID:   uuid.New(),
		Name: requestedName,
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

	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return err
	}

	err = s.db.ResetFeedFollows(context.Background())
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

func handlerAggregator(s *state, cmd command) error {
	fetchedFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v", *fetchedFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 2 {
		fmt.Println("2 arguments required")
		os.Exit(1)
	}

	feedNameInput, urlInput := cmd.args[0], cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   feedNameInput,
		Url:    urlInput,
		UserID: currentUser.ID,
	}

	createdFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParam := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: currentUser.ID,
		FeedID: createdFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParam)
	if err != nil {
		return err
	}

	fmt.Println("feed created. Feed data:")
	fmt.Printf("name: %s\n", createdFeed.Name)
	fmt.Printf("url: %s\n", createdFeed.Url)
	fmt.Printf("user id: %s\n", createdFeed.UserID)

	os.Exit(0)
	return nil
}

func handlerFeedsList(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		// looks costly
		creator, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s | Url: %s | Creator: %s\n", feed.Name, feed.Url, creator.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, creator database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("a single argument needed")
		return fmt.Errorf("a single argument needed")
	}
	inputUrl := cmd.args[0]

	feed, err := s.db.GetFeedFromUrl(context.Background(), inputUrl)
	if err != nil {
		return err
	}

	feedFollowParam := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: creator.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParam)
	if err != nil {
		return err
	}

	fmt.Println("Following success! data:")
	fmt.Printf("feed name: %s | current user: %s", feed.Name, s.config.CurrentUsername)
	return nil
}

func handlerFollowsList(s *state, cmd command, currentUser database.User) error {
	followingFeeds, err := s.db.GetFeedFollowForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	for _, feed := range followingFeeds {
		fmt.Printf("feed name: %s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Feed Url required")
		return fmt.Errorf("Feed Url required")
	}

	url := cmd.args[0]

	targetFeed, err := s.db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return err
	}

	deleteParam := database.DeleteFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: targetFeed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteParam)
	if err != nil {
		return err
	}

	fmt.Println("Unfollow successful")
	return nil
}
