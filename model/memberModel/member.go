package memberModel

type MemberEdit struct {
	Uid     string `json:"uid" gorm:"type:varchar(255);not null"`
	VipTime string `json:"vip_time" gorm:"type:varchar(255);not null"`
}
