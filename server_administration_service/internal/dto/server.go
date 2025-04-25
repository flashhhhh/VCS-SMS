package dto

type ServerAddress struct {
	ID int `json:"id" gorm:"primary_key"`
	IPv4      string `json:"ip_address" gorm:"not null"`
	Port      int    `json:"port" gorm:"not null"`
}

type ServerFilter struct {
	ServerID string `json:"server_id"`
	ServerName string `json:"server_name"`
	Status	 string `json:"status"`
	IPv4	  string `json:"ipv4"`
	Port	  int    `json:"port"`
}