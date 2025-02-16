# Gator

Gator is a RSS command line tool for

1. Fetching feeds
2. Storing feeds
3. Following feeds

## How to install

1. Have these installed

   - [go](https://go.dev/doc/install) (1.24+)
   - [postgres](https://www.postgresql.org/download/) (16.6+)
   - [goose](https://github.com/pressly/goose), using `go install github.com/pressly/goose/v3/cmd/goose@latest`

2. Clone this repo `git clone https://github.com/WaronLimsakul/gator.git`
3. Set up a postgres database
   - in the local machine, create new postgres database in any port you want
4. Set up the config file
   - Go will look for a config file in `~/.gatorconfig.json`
   - Users can setup a username by running `gator register [username]`
   - Users can setup a database url by command `gator setdb [dburl]`
5. In the root of the program, run `goose -dir ./sql/schema/ up`
6. In the root, Install gator using `go install .`

## Commands

### Account-related command

- Register with username: `register [username]`
- After registering, user can login: `login [username]`
- List registered users: `users`

### Feeds commands

- Reset entire database (rows, not schema): `reset`
- Scrape the followed feeds and save posts in database every specified time: `agg [time(opt)]`
  - time can be in format ?s ?m ?h
- Add feed by url: `addfeed [feed name] [url]`
- List all feeds: `feeds`
- Follow an available feed in database: `follow [url]`
- Display following feeds: `following`
- Unfollow feed: `unfollow [url]`
- Display stored posts: `browse`
