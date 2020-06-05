package models

import (
	"time"
)

//BaseModel using own model instead of gorm.Model
type BaseModel struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Node struct {
	BaseModel

	NodeId       string `gorm:"type:varchar(100);unique_index;not null" json:"node_id"`
	ContainerId  string `gorm:"type:varchar(100);unique_index;not null" json:"container_id"`
	NodeFolder   string `gorm:"type:varchar(100);not null" json:"node_folder"`
	NodeManifest string `gorm:"type:varchar(100);not null" json:"node_manifest"`
}
