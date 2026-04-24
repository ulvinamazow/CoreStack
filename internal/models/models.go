package models

import (
	"time"
)

type User struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username          string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Name              string     `gorm:"type:varchar(100);not null" json:"name"`
	Gmail             string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"gmail"`
	PasswordHash      string     `gorm:"type:varchar(255);not null" json:"-"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified"`
	VerificationToken *string    `gorm:"type:varchar(255)" json:"-"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	IsAdmin           bool       `gorm:"default:false" json:"is_admin"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	TokenHash string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Category struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	ParentID  *uint      `json:"parent_id,omitempty"`
	Parent    *Category  `gorm:"foreignKey:ParentID" json:"-"`
	Children  []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Product struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID     uint      `gorm:"not null" json:"seller_id"`
	Seller       User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	CategoryID   uint      `gorm:"not null" json:"category_id"`
	Category     Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name         string    `gorm:"type:varchar(200);not null" json:"name"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Discount     int       `gorm:"default:0" json:"discount"`
	Stock        int       `gorm:"default:0" json:"stock"`
	Description  string    `gorm:"type:text" json:"description"`
	AverageRating float64  `gorm:"-" json:"average_rating,omitempty"`
	ReviewCount  int       `gorm:"-" json:"review_count,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Product) DiscountedPrice() float64 {
	if p.Discount <= 0 {
		return p.Price
	}
	return p.Price * float64(100-p.Discount) / 100
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null;check:quantity > 0" json:"quantity"`
	AddedAt   time.Time `json:"added_at"`
}

type Order struct {
	ID                    uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                uint        `gorm:"not null" json:"user_id"`
	User                  User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalAmount           float64     `gorm:"type:decimal(10,2)" json:"total_amount"`
	StripePaymentIntentID *string     `gorm:"type:varchar(255)" json:"stripe_payment_intent_id,omitempty"`
	Status                string      `gorm:"type:varchar(20);default:pending" json:"status"`
	Items                 []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt             time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID              uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID         uint    `gorm:"not null" json:"order_id"`
	ProductID       uint    `gorm:"not null" json:"product_id"`
	Product         Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity        int     `gorm:"not null" json:"quantity"`
	UnitPrice       float64 `gorm:"type:decimal(10,2)" json:"unit_price"`
	DiscountApplied int     `json:"discount_applied"`
}

type Review struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	OrderID   uint      `gorm:"not null" json:"order_id"`
	Rating    int       `gorm:"not null;check:rating between 1 and 5" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
