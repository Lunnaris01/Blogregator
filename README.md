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
Install the `gator` CLI using the `go install` command:

```bash
go install github.com/Lunnaris01/Blogregator/cmd/gator@latest
```
## Usage
Once set up, you can run Gator via the command line. Below are the available commands and their descriptions:

### Commands Overview

1. **`gator login`**
   - Logs a user into the application. This command prompts for credentials.
   - Example:  
     ```bash
     gator login
     ```

2. **`gator register`**
   - Registers a new user with the system.
   - Example:  
     ```bash
     gator register
     ```

3. **`gator reset`**
   - Resets your application's settings or database to its initial state.
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
   - A command to perform some aggregate function (be sure to elaborate what `agg` does in your project!).
   - Example:  
     ```bash
     gator agg
     ```

6. **`gator addfeed`**
   - Adds a new feed for the logged-in user. This command requires the user to be logged in.
   - Example:  
     ```bash
     gator addfeed
     ```

7. **`gator feeds`**
   - Lists all feeds in the system.
   - Example:  
     ```bash
     gator feeds
     ```

8. **`gator follow`**
   - Allows a logged-in user to follow a particular feed.
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
      gator explore 20
      ```

12. **`gator posts`**
    - Lists every single post saved (mostly for debug purposes).
    - Example:  
      ```bash
      gator posts
      ```

