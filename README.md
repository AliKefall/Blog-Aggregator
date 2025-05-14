# 🐊 Gator – Blog Aggregator CLI

Gator is a terminal-based blog RSS feed aggregator and PostgreSQL explorer. It allows developers to interact with their database and RSS feeds directly from the command line, without needing to use a GUI like pgAdmin.

---

## 🚀 Features

* **Database Exploration:** Browse tables, schemas, and data directly from your terminal.
* **Simplified Querying:** Run essential queries without writing raw SQL.
* **CLI-First Workflow:** Designed to integrate into a developer's terminal-based workflow.
* **Easy Configuration:** Connection settings handled through a JSON config file.
* **Feed Management:** Subscribe to, unfollow, and view RSS feeds directly through the CLI.

---

## 🛠 Requirements

* Go (Golang)
* PostgreSQL

---

## ⚙️ Environment Setup

1. Create a `.gatorconfig.json` file in your home directory.
2. Add the following content with your own connection information:

```json
{
  "current_user_name": "your_username",
  "db_url": "postgres://youruser:yourpass@localhost:5432/yourdb"
}
```

> Why not automate it? Because exposing database credentials automatically is a security risk. Users must provide their own.

---

## 🪡 Optional: Database Migration with Goose

### Install Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Run Migrations:

```bash
goose -dir ./sql/schema postgres "your_connection_string" up
```

### Rollback Migrations:

```bash
goose -dir ./sql/schema postgres "your_connection_string" down
```

---

## 📁 Config File Location

| OS          | Path                                      |
| ----------- | ----------------------------------------- |
| Windows     | `C:\Users\YourUsername\.gatorconfig.json` |
| Linux/macOS | `~/.gatorconfig.json`                     |

---

## 📖 CLI Commands

| Command                | Description                                              |
| ---------------------- | -------------------------------------------------------- |
| `login <username>`     | Log in with a registered username.                       |
| `register <username>`  | Register a new user in the database.                     |
| `users`                | List all registered users.                               |
| `reset`                | Clears all data from tables but retains structure.       |
| `addfeed <name> <url>` | Add a new RSS feed by name and URL.                      |
| `feeds`                | List all available feeds.                                |
| `follow` / `unfollow`  | Follow or unfollow a feed.                               |
| `following`            | Show currently followed feeds.                           |
| `agg <interval>`       | Aggregate posts at defined intervals (e.g., `1m`, `1h`). |
| `browse <limit>`       | View fetched posts (default limit: 2).                   |

---

## 🎯 Example Usage

```bash
register johndoe
login johndoe
addfeed golangweekly https://golangweekly.com/rss
follow https://golangweekly.com/rss
agg 1m
browse 5
```

---

## 📃 Notes

* All commands are run from the terminal.
* Goose is optional but recommended for managing schema changes.
* The `.gatorconfig.json` file is required to connect to your database.

---

Gator is designed for developers who prefer working with data without leaving their terminal.

