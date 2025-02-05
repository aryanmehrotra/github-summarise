# GitHub Activity Fetcher

This project fetches various GitHub activities (issues created, pull requests reviewed, pull requests merged, commits created, and comments) for a specific user on a given date. The data is fetched using GitHub's official API and served through a REST API built with GoFr.

## Features

- Fetch GitHub activities for a given user and date:
    - Issues created
    - Pull requests reviewed
    - Pull requests merged
    - Commits created
    - Comments made
- Each activity includes the title and a link to the relevant GitHub resource (issue, PR, commit, etc.).

## Requirements

- Go 1.18+
- GitHub personal access token (with `repo` and `user` permissions)
- GoFr package

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/github-activity-fetcher.git
   cd github-activity-fetcher
   ```
   
2. Install dependencies
   ```
    go mod tidy
    ```

3. Run the application
   ```
    go run main.go
    ```