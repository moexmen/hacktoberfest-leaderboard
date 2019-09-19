# Hacktoberfest Leaderboard

Get the PR counts and avatar URLs for a list of authors defined in the `AUTHORS` environment variable and return as JSON.

## Building and Running

Install [Docker](https://docs.docker.com/install/) and [Docker Compose](https://docs.docker.com/compose/install/)

Run `docker-compose up`. Use `docker-compose up --build` if the image needs to be rebuilt.

## Environment Variables

Copy `.env.sample` to `.env`.

To increase GitHub's rate limit and allow up to 30 authors on the leaderboard, obtain a [personal access token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) from GitHub.
Without the token, only 10 authors can be added as the [search API's rate limit](https://developer.github.com/v3/search/#rate-limit) will be reached.

Assign the personal access token to the `GHTOKEN` variable.
Add the authors to be included on the leaderboard to the `AUTHORS` variable, separated by the colon character `:`.

## Retrieving Data

With the Docker container running, open a browser and go to [https://localhost:4000/leaderboard](https://localhost:4000/leaderboard).
