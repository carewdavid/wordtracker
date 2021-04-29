#! /usr/bin/python
import time
import os
import sqlite3
dbpath = os.environ.get("WORDCOUNT_DB", "wordcound.db")
db = sqlite3.connect(dbpath)
cur = db.cursor()

#There can be multiple entries for one date, so we can't just use AVG.
#Actually this isn't right either, but it'll take fancier date calculations to fix it
#so I'm just going to put it on the todo list
cur.execute("SELECT sum(words), SUM(words) / COUNT(DISTINCT date) from wordcount")
total = cur.fetchone()
month = time.strftime("%Y-%m")
cur.execute('SELECT sum(words), SUM(words) / COUNT(DISTINCT date) from wordcount WHERE strftime("%Y-%m", date) = ?', [month])
monthly = cur.fetchone()
if monthly is None:
    monthly = (0, 0)
week = time.strftime("%Y-%W")
cur.execute('SELECT sum(words), SUM(words) / COUNT(DISTINCT date) from wordcount WHERE strftime("%Y-%W", date) = ?', [week])
weekly = cur.fetchone()
if weekly is None:
    weekly = (0, 0)

print(f'Total to date: {total[0]} words. Average {total[1]:.1f} / day')
print(f'Total this month: {monthly[0]}. {monthly[1]:.1f} / day')
print(f'Total this week: {weekly[0]}. {weekly[1]:.1f} / day')
db.close()

