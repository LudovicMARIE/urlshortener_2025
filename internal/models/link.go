package models

import "time"

type Link struct {
	ID        uint      `gorm:"primaryKey"`
	Shortcode string    `gorm:"uniqueIndex;size:10"`
	LongURL   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
