import json
from pathlib import Path
import sqlite3

source = "data.json"
dest = "db.sqlite"

data = json.loads(Path(source).read_text())
templates = data["templates"]
configs = data["configs"]

conn = sqlite3.connect(dest)

cur = conn.cursor()

for template in templates:
    query = "INSERT INTO timestamps (first, rest) VALUES (?, ?)"
    first, rest = template
    cur.execute(query, (first, rest))

for key, config in configs.items():
    query = "INSERT INTO courses VALUES (?, ?, ?, ?)"
    print(key)
    cur.execute(query, (key, config["slug"], config["template"], "codewithmosh"))
    for i, url in enumerate(config["magnets"], start=1):
        category = "magnet" if url.startswith("magnet") else "gdrive"
        query = "INSERT INTO course_urls (course_id, url, category, position) VALUES (?, ?, ?, ?)"
        cur.execute(query, (key, url, category, i))
    
conn.commit()
conn.close()