package models

import "time"

type Column struct {
	Id        int64     `json:"id"`
	TableId   int64     `sql:"not null" json:"table_id"`
	TableName string    `sql:"size:45;not null" json:"table_name"`
	Name      string    `sql:"size:45;not null" json:"name"`
	Alias     string    `sql:"size:45" json:"alias"`
	Type      string    `sql:"size:20" json:"type"`
	Size      int       `json:"size"`
	Filter    string    `sql:"size:127" json:"filter"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
}
