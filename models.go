package undersky

import "time"

// User 는 회원 모델입니다.
type User struct {
	ID int64 `gorm:"primary_key"`

	UUID     string `gorm:"type:varchar(36);not null;unique_index"`
	Email    string `gorm:"type:varchar(255);not null;unique_index"`
	Username string `gorm:"type:varchar(40);not null;unique_index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCredential 함수는 유저에게 종속되는 새로운 인증 정보를 생성합니다.
func (u *User) NewCredential() *Credential {
	return &Credential{
		UserID:      u.ID,
		AccessToken: generateRandomString(40),
		SecretToken: generateRandomString(40),
		ValidUntil:  time.Now().AddDate(0, 0, 14).UTC(),
	}
}

// Credential 은 회원의 로그인 토큰 모델입니다.
type Credential struct {
	ID int64 `gorm:"primary_key"`

	UserID      int64     `gorm:"index;not null"`
	AccessToken string    `gorm:"type:varchar(255);not null;unique_index"`
	SecretToken string    `gorm:"type:varchar(255);not null"`
	ValidUntil  time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// HasExpired 함수는 주어진 인증 정보의 만료 여부를 반환합니다.
func (c *Credential) HasExpired() bool {
	return c.ValidUntil.Before(time.Now())
}
