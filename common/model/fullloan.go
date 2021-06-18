package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LeadFullLoan struct {
	gorm.Model
	LeadID      int            `gorm:"column:lead_id;type:int(9);not null" json:"lead_id"`
	RequestID   string         `gorm:"column:request_id;type:varchar(20);not null;unique" json:"request_id"`
	PartnerCode string         `gorm:"column:partner_code;type:varchar(20);not null" json:"partner_code"`
	Document    datatypes.JSON `gorm:"column:document;type:longtext" json:"document"`
}

// TableName sets the insert table name for this struct type
func (l *LeadFullLoan) TableName() string {
	return "vicidial_list_full_loan"
}
