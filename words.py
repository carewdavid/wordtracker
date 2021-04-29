#! /usr/bin/python
import sqlite3
import sys
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

db = sqlite3.connect("wordcount.db")
cur = db.cursor()
cur.execute("INSERT INTO wordcount VALUES(?, ?, ?)", (date.isoformat(), words, desc))
db.commit()
db.close()

