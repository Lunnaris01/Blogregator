# Gator CLI

Gator is a CLI tool designed to Aggregate Posts from RSS feeds. Built with **Go** and backed by the power of **Postgres**, Gator showcases efficient backend development and CLI tooling.

## Technologies
- [Go](https://go.dev): A statically-typed, compiled language for modern applications.
- [Postgres](https://www.postgresql.org): A robust, open-source relational database system.

## Requirements
To run this program, you'll need:
- **Postgres** installed. You can find installation details [here](https://www.postgresql.org/download/).
- **Go** installed. For installation instructions, visit [the official Go documentation](https://go.dev/doc/install).

## Installation
Install the `gator` CLI using the `go install` command (might not work if you dont use linux).:

```bash
go install github.com/Lunnaris01/Blogregator/cmd/gator@latest
```
## Usage
Once set up, you can run Gator via the command line. Below are the available commands and their descriptions:

### Commands Overview

1. **`gator login`**
   - Logs a user into the application.
   - Example:  
     ```bash
     gator login <username>
     ```

2. **`gator register`**
   - Registers a new user with the system.
   - Example:  
     ```bash
     gator register <username>
     ```

3. **`gator reset`**
   - Resets your application's settings or database to its initial state. Use with caution!
   - Example:  
     ```bash
     gator reset
     ```

4. **`gator users`**
   - Prints a list of all registered users in the system.
   - Example:  
     ```bash
     gator users
     ```

5. **`gator agg`**
   - Agg fetches the feed which hasn't been fetched for the longest time and saves all posts from the feed to the database. Note that only feeds the logged in user is following are considered.
   - Example:  
     ```bash
     gator agg <interval>
     ```

6. **`gator addfeed`**
   - Adds a new feed by url for the logged-in user. 
   - Example:  
     ```bash
     gator addfeed <feed_url>
     ```

7. **`gator feeds`**
   - Lists all feeds in the system.
   - Example:  
     ```bash
     gator feeds
     ```

8. **`gator follow`**
   - Allows a logged-in user to follow a particular feed. This is automatically called for the current user if he is adding a new feed.
   - Example:  
     ```bash
     gator follow <feed_name>
     ```

9. **`gator following`**
   - Displays a list of feeds the currently logged-in user is following.
   - Example:  
     ```bash
     gator following
     ```

10. **`gator unfollow`**
    - Allows a logged-in user to unfollow a feed.
    - Example:  
      ```bash
      gator unfollow <feed_name>
      ```

11. **`gator explore`**
    - Allows the logged in user to list his N most recently published Posts from the followed Feeds.
    - Example:  
      ```bash
      gator explore <N>
      ```

12. **`gator posts`**
    - Lists every single post saved (mostly for debug purposes).
    - Example:  
      ```bash
      gator posts
      ```

