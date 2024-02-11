package domain

import "time"

type Comment struct {
	Id      int64  `json:"id"`
	User    User   `json:"user"`
	BizId   int64  `json:"biz_id,omitempty"`
	Biz     string `json:"biz,omitempty"`
	Content string `json:"content,omitempty"`

	RootComment   *Comment   `json:"root_comment,omitempty"`
	ParentComment *Comment   `json:"parent_comment,omitempty"`
	Children      []*Comment `json:"children,omitempty"`

	Utime time.Time `json:"utime"`
	Ctime time.Time `json:"ctime"`
}

type User struct {
	Id   int64
	Name string
}
