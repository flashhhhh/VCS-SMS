package dto

type ServerAddress struct {
	ServerID string `json:"server_id" gorm:"primary_key"`
	IPv4      string `json:"ip_address" gorm:"not null"`
	Port      int    `json:"port" gorm:"not null"`
}