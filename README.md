A simple tracker for writing productivity.

Disclaimer: you should probably just use a spreadsheet.

Set the environment variable `WORDCOUNT_DB` to the path of where you want to keep your database.

Log your writing with words.py. It's usage is:
```
$ words.py [-y] n ["description"]
```
The only required argument is n, the number of words to log. By default, it will be logged under the current date, if `-y` is present, the previous day will be used as the date instead. A description, e.g. if you have multiple projects can optionally be added as well. The words in the description are not counted in the log.

report.py allows you to view your data. It displays the all time total count, as well as counts for the past calendar month and week.