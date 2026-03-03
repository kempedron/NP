package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string       `gorm:"type:varchar(50);not null;unique" json:"username"`
	PasswordHash string       `gorm:"type:varchar(100);not null" json:"-"`
	BankAccount  *BankAccount `json:"bank_account,omitempty"`
	Cart         *Cart        `json:"cart,omitempty"`
}

func (u *User) HashPassword(password string) error {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedpassword)
	return nil
}

func (u User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (User) TableName() string {
	return "users"
}

type Cart struct {
	gorm.Model
	UserID    uint        `gorm:"not null;unique" json:"user_id"`
	User      *User       `json:"-"`
	Items     []*CartItem `json:"items"`
	TotalCost uint        `gorm:"type:BIGINT;not null;default:0" json:"total_cost"`
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	gorm.Model
	CartID    uint     `gorm:"not null;uniqueIndex:idx_cart_product" json:"cart_id"`
	ProductID uint     `gorm:"not null;uniqueIndex:idx_cart_product" json:"product_id"`
	Quantity  uint     `gorm:"type:SMALLINT;not null;default:1" json:"quantity"`
	Cart      *Cart    `json:"-"`
	Product   *Product `json:"product,omitempty"`
}

func (CartItem) TableName() string {
	return "carts_item"
}

type Product struct {
	gorm.Model
	Name        string      `gorm:"type:varchar(50);not null" json:"name"`
	Price       uint        `gorm:"type:BIGINT;not null" json:"price"`
	Description string      `gorm:"type:varchar(100);not null" json:"description"`
	CartItems   []*CartItem `json:"-"`
}

func (Product) TableName() string {
	return "products"
}

type BankAccount struct {
	gorm.Model
	Balance uint64 `gorm:"type:BIGINT;not null" json:"balance"`
	UserID  uint   `gorm:"not null;unique" json:"user_id"`
	User    *User  `json:"-"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}

type Donate struct {
	gorm.Model
	Username  string `gorm:"type:varchar(50);not null" json:"username"`
	MoneySumm uint   `gorm:"type:BIGINT;not null" json:"price"`
	Category  string `gorm:"type:varchar(100);not null" json:"category"`
}

func (Donate) TableName() string {
	return "donates"
}
