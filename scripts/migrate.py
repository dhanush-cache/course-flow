import json
from pathlib import Path
import sqlite3

base_path = Path(__file__).resolve().parent.parent
source = base_path / "scripts" / "data.json"
dest = Path.home() / ".local" / "share" / "course-flow" / "db.sqlite"
dest.parent.mkdir(parents=True, exist_ok=True)

data = json.loads(source.read_text())
templates = data["templates"]
configs = data["configs"]

conn = sqlite3.connect(dest)
cur = conn.cursor()

# Enable foreign keys
cur.execute("PRAGMA foreign_keys = ON;")

# Populate platforms
platforms = [
    ("codewithmosh", "Code With Mosh", "https://codewithmosh.com"),
    ("dreamsofcode", "Dreams of Code", "https://dreamsofcode.io")
]
cur.executemany("INSERT OR IGNORE INTO platforms (id, title, url) VALUES (?, ?, ?)", platforms)

for template in templates:
    query = "INSERT INTO timestamps (first, rest) VALUES (?, ?)"
    first, rest = template
    cur.execute(query, (first, rest))

for key, config in configs.items():
    platform_id = "dreamsofcode" if key == "go-cli" else "codewithmosh"
    # template index to timestamp_id (index + 1)
    timestamp_id = config["template"] + 1
    query = "INSERT INTO courses VALUES (?, ?, ?, ?)"
    print(key)
    cur.execute(query, (key, config["slug"], timestamp_id, platform_id))
    for i, url in enumerate(config["magnets"], start=1):
        category = "magnet" if url.startswith("magnet") else "gdrive"
        query = "INSERT INTO course_urls (course_id, url, category, position) VALUES (?, ?, ?, ?)"
        cur.execute(query, (key, url, category, i))
    
conn.commit()
conn.close()