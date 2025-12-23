package userModel

type Wallet struct {
	Income     float64 `gorm:"type:decimal(10,2);not null;default:0" json:"income"`
	Balance    float64 `gorm:"type:decimal(10,2);not null;default:0" json:"balance"`
	Withdrawal float64 `gorm:"type:decimal(10,2);not null;default:0" json:"withdrawal"`
}
