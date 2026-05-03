# Course-Flow

Course-Flow is a command-line interface (CLI) tool designed to help you efficiently manage your collection of educational courses.

## Setup Instructions

### Prerequisites

- Go (1.21 or later)
- Python 3
- [golang-migrate](https://github.com/golang-migrate/migrate)
- SQLite3

### Installation

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd course-flow
   ```

2. **Install Go dependencies:**

   ```bash
   go mod download
   ```

3. **Apply Database Migrations:**
   Ensure you have `migrate` installed and run:

   ```bash
   migrate -path db/migrations -database "sqlite3://$HOME/.local/share/course-flow/db.sqlite" up
   ```

4. **Populate Database (Optional):**
   Run the migration script to populate the database with initial data:
   ```bash
   python3 scripts/migrate.py
   ```

### Usage

Run the application:

```bash
go run main.go
```
