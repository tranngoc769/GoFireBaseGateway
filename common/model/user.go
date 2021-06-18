package model

import "time"

type User struct {
	Id          string    `json:"id" gorm:"primaryKey;column:id;type:int(11) AUTO_INCREMENT;not null"`
	Username    string    `json:"username" gorm:"unique;column:username;type:varchar(20);not null"`
	Level       string    `json:"level" gorm:"column:level;type:int(4);not null"`
	Password    string    `json:"password" gorm:"column:password;type:varchar(200);not null"`
	ApiKey      string    `json:"api_key" gorm:"column:api_key;type:varchar(200);null"`
	CreateAt    time.Time `json:"created_at" gorm:"column:created_at;type:TIMESTAMP"`
	PartnerCode string    `json:"partner_code" gorm:"column:partner_code;type:varchar(20)"`
}

func (User) TableName() string {
	return "vicidial_api_user"
}

type UserAuthRes struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	ApiKey      string `json:"api_key"`
	Level       string `json:"level"`
	PartnerCode string `json:"partner_code"`
}

type AccessToken struct {
	ClienID      string    `json:"client_id"`
	UserID       string    `json:"user_id"`
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Createdtime  time.Time `json:"create_at"`
	Expiredtime  int       `json:"expire_at"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
	JWT          string    `json:"jwt"`
}
