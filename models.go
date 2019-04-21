package undersky

import "time"

// User 는 회원 모델입니다.
type User struct {
	ID    int64  `gorm:"primary_key"`
	UUID  string `gorm:"type:varchar(36);unique_index"`
	Email string `gorm:"type:varchar(255);unique_index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
