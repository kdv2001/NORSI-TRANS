package models

type Notion struct {
	Id          int64  `json:"notionId" db:"notionId"`
	UserId      int64  `json:"userId" db:"userId"`
	Information string `json:"notion" db:"notion"`
}
