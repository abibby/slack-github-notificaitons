package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var ghClient *github.Client

func getGHClient(ctx context.Context) *github.Client {
	if ghClient == nil {
		token := os.Getenv("GITHUB_PAT")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		ghClient = github.NewClient(tc)
	}
	return ghClient
}

func GitHubNotifications(since time.Time) ([]*github.Notification, error) {
	ctx := context.Background()
	client := getGHClient(ctx)

	notifs, _, err := client.Activity.ListNotifications(ctx, &github.NotificationListOptions{
		Since: since,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifs, err
}
