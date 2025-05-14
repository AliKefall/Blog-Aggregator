package main

import (
	"context"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github/AliKefall/gator/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd command) error {
	limit := 2
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("Invalid Limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Couldn't get posts for user: %w", err)
	}
	p := bluemonday.StripTagsPolicy()

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)

		cleanDescription := p.Sanitize(post.Description.String)
		fmt.Printf("    %v\n", cleanDescription)

		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
