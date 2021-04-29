#! /usr/bin/python
import time
import sqlite3

db = sqlite3.connect("wordcount.db")
cur = db.cursor()

cur.execute("SELECT sum(words) from wordcount")
total = cur.fetchone()[0]
month = time.strftime("%Y-%m")
cur.execute('SELECT sum(words) from wordcount WHERE strftime("%Y-%m", date) = ?', [month])
try:
    monthly = cur.fetchone()[0]
except TypeError:
    monthly = None
week = time.strftime("%Y-%W")
cur.execute('SELECT sum(words) from wordcount WHERE strftime("%Y-%W", date) = ?', [week])
try:
    weekly = cur.fetchone()[0]
except TypeError:
    weekly = None

print(f'Total to date: {total}')
print(f'Total this month: {monthly}')
print(f'Total this week: {weekly}')
db.close()

