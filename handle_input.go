package main

import (
	"context"
	"fmt"
	"github/AliKefall/gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("User created succesfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("User switched succesfully:")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Arguments of the reset command must be empty!")
	}

	err := s.db.CleanTable(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database is cleaned!")
	return nil

}

func handlerList(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Arguments of the list command must be empty")
	}

	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {

		if user.Name == s.cfg.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}
func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Arguments of the feeds command must be empty!")
	}
	feeds, err := s.db.ListFeedsWithUsers(context.Background())

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s\n", feed.FeedName)
		fmt.Printf("URL: %s\n", feed.FeedUrl)
		fmt.Printf("Created By: %s\n\n", feed.UserName)
	}
	return nil
}

func handlerFetchedData(s *state, cmd command) error {
	fetchedData, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", fetchedData)
	return nil
}

func handlerFollow(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <feed_url>", cmd.Name)

	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("Feed follow created")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil

}

func handlerUnfollow(s *state, cmd command) error {
	userId, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: userId.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return err
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
