package record

type Record struct {
	Date  int64  `json:date`
	Words int    `json:date`
	Desc  string `json:desc`
}
