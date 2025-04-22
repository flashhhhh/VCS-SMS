package domain

import "time"

type Server struct {
	ServerID string `json:"server_id" gorm:"primary_key"`
	ServerName string `json:"server_name" gorm:"not null"`
	Status string `json:"status" gorm:"not null"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	LastUpdated time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	IPAddress string `json:"ip_address" gorm:"not null"`
	Port int `json:"port" gorm:"not null"`
}