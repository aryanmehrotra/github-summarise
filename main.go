package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v59/github"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2"
)

const (
	dateFormat = "02-01-2006" // DD-MM-YYYY
)

// Response struct to encapsulate all activities
type Response struct {
	IssuesCreated  []Activity `json:"issues_created"`
	PRsReviewed    []Activity `json:"prs_reviewed"`
	PRsMerged      []Activity `json:"prs_merged"`
	CommitsCreated []Activity `json:"commits_created"`
	Comments       []Activity `json:"comments"`
}

// Activity struct to represent a single activity (issue, PR, commit, or comment)
type Activity struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func main() {
	app := gofr.New()

	githubToken := app.Config.Get("GITHUB_TOKEN")
	githubUser := app.Config.Get("GITHUB_USER")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	app.GET("/fetch-activities", func(c *gofr.Context) (interface{}, error) {
		dateStr := c.Param("date")
		date, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}

		// Channels to collect results
		issuesCh := make(chan []Activity)
		prsReviewedCh := make(chan []Activity)
		prsMergedCh := make(chan []Activity)
		commitsCh := make(chan []Activity)
		commentsCh := make(chan []Activity)

		go func() {
			issuesCh <- fetchIssuesCreated(ctx, client, githubUser, date)
		}()
		go func() {
			prsReviewedCh <- fetchPRsReviewed(ctx, client, githubUser, date)
		}()
		go func() {
			prsMergedCh <- fetchPRsMerged(ctx, client, githubUser, date)
		}()
		go func() {
			commitsCh <- fetchCommitsCreated(ctx, client, githubUser, date)
		}()
		go func() {
			commentsCh <- fetchComments(ctx, client, githubUser, date)
		}()

		// Collect results from channels
		response := Response{
			IssuesCreated:  <-issuesCh,
			PRsReviewed:    <-prsReviewedCh,
			PRsMerged:      <-prsMergedCh,
			CommitsCreated: <-commitsCh,
			Comments:       <-commentsCh,
		}

		// Close channels (good practice after receiving from them)
		close(issuesCh)
		close(prsReviewedCh)
		close(prsMergedCh)
		close(commitsCh)
		close(commentsCh)

		return response, nil
	})

	app.Run()
}

func fetchIssuesCreated(ctx context.Context, client *github.Client, user string, date time.Time) []Activity {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	query := fmt.Sprintf("author:%s type:issue created:%s", user, date.Format("2006-01-02"))
	results, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		fmt.Printf("Error fetching issues: %v\n", err)
		return nil
	}

	var activities []Activity
	for _, issue := range results.Issues {
		activities = append(activities, Activity{
			Title: issue.GetTitle(),
			URL:   issue.GetHTMLURL(),
		})
	}

	return activities
}

func fetchPRsReviewed(ctx context.Context, client *github.Client, user string, date time.Time) []Activity {
	opts := &github.SearchOptions{Sort: "updated", Order: "desc"}
	query := fmt.Sprintf("reviewed-by:%s type:pr updated:%s", user, date.Format("2006-01-02"))
	results, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		fmt.Printf("Error fetching PRs reviewed: %v\n", err)
		return nil
	}

	var activities []Activity
	for _, pr := range results.Issues {
		activities = append(activities, Activity{
			Title: pr.GetTitle(),
			URL:   pr.GetHTMLURL(),
		})
	}

	return activities
}

func fetchPRsMerged(ctx context.Context, client *github.Client, user string, date time.Time) []Activity {
	opts := &github.SearchOptions{Sort: "updated", Order: "desc"}
	query := fmt.Sprintf("author:%s type:pr is:merged updated:%s", user, date.Format("2006-01-02"))
	results, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		fmt.Printf("Error fetching merged PRs: %v\n", err)
		return nil
	}

	var activities []Activity
	for _, pr := range results.Issues {
		activities = append(activities, Activity{
			Title: pr.GetTitle(),
			URL:   pr.GetHTMLURL(),
		})
	}

	return activities
}

func fetchCommitsCreated(ctx context.Context, client *github.Client, user string, date time.Time) []Activity {
	repos, _, err := client.Repositories.List(ctx, user, nil)
	if err != nil {
		fmt.Printf("Error fetching repositories: %v\n", err)
		return nil
	}

	var activities []Activity
	for _, repo := range repos {
		commits, _, err := client.Repositories.ListCommits(ctx, repo.GetOwner().GetLogin(), repo.GetName(), &github.CommitsListOptions{
			Author: user,
			Since:  date,
			Until:  date.Add(24 * time.Hour),
		})

		if err != nil {
			continue
		}

		for _, commit := range commits {
			activities = append(activities, Activity{
				Title: commit.GetCommit().GetMessage(),
				URL:   commit.GetHTMLURL(),
			})
		}
	}

	return activities
}

func fetchComments(ctx context.Context, client *github.Client, user string, date time.Time) []Activity {
	opts := &github.SearchOptions{Sort: "updated", Order: "desc"}
	query := fmt.Sprintf("commenter:%s updated:%s", user, date.Format("2006-01-02"))
	results, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		fmt.Printf("Error fetching comments: %v\n", err)
		return nil
	}

	var activities []Activity
	for _, issue := range results.Issues {
		activities = append(activities, Activity{
			Title: issue.GetTitle(),
			URL:   issue.GetHTMLURL(),
		})
	}

	return activities
}
