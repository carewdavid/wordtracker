#! /usr/bin/python
import sqlite3
import sys
import os
import time
from datetime import date, timedelta

date = date.today()
if sys.argv[1] == '-y':
    date = date - timedelta(days=1)
    del sys.argv[1]

try:
    words = int(sys.argv[1])
except:
    print("First argument must be a number")
    sys.exit(1)

desc = None
if len(sys.argv) > 2:
    desc = sys.argv[2]
dbpath = os.environ.get("WORDCOUNT_DB", "wordcount.db")
db = sqlite3.connect(dbpath)
cur = db.cursor()
cur.execute("CREATE TABLE IF NOT EXISTS wordcount(date STRING NOT NULL, words INT NOT NULL, desc STRING)")
cur.execute("INSERT INTO wordcount VALUES(?, ?, ?)", (date.isoformat(), words, desc))
db.commit()
db.close()

