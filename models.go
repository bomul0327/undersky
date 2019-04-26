package undersky

import (
	"time"

	"github.com/google/uuid"
	"github.com/looplab/fsm"
)

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

// Game 은 게임 모델입니다.
type Game struct {
	ID int64 `gorm:"primary_key"`

	GameID      string `gorm:"type:varchar(10);not null;unique_index"`
	Title       string `gorm:"type:varchar(64);not null"`
	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Match 는 두 유저간의 경기 모델입니다.
type Match struct {
	ID   int64  `gorm:"primary_key"`
	UUID string `gorm:"type:varchar(36);not null;unique_index"`

	// submitted | initializing | started | finished | failed
	State *fsm.FSM `gorm:"type:varchar(16);not null"`

	GameID                 string `gorm:"type:varchar(10);not null;index"`
	MatchUUID              string `gorm:"type:varchar(36);not null;unique_index"`
	PlayerID               int64  `gorm:"not null;index"`
	PlayerSubmissionID     int64  `gorm:"not null;index"`
	CompetitorID           int64  `gorm:"not null;index"`
	CompetitorSubmissionID int64  `gorm:"not null"`

	TotalRound    int
	PlayerWin     int
	CompetitorWin int

	// win | lose | draw | error
	Result       string `gorm:"type:varchar(16)"`
	ErrorMessage string `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMatch 는 두 유저간의 경기를 만듭니다.
func NewMatch(gameID string, playerID, playerSubID, compID, compSubID int64) *Match {
	u, _ := uuid.NewRandom()
	return &Match{
		UUID: u.String(),

		GameID:                 gameID,
		PlayerID:               playerID,
		PlayerSubmissionID:     playerSubID,
		CompetitorID:           compID,
		CompetitorSubmissionID: compSubID,

		State: fsm.NewFSM(
			"submitted",
			fsm.Events{
				{Name: "init", Src: []string{"submitted"}, Dst: "initializing"},
				{Name: "start", Src: []string{"initializing"}, Dst: "started"},
				{Name: "finish", Src: []string{"started"}, Dst: "finished"},
				{Name: "fail", Src: []string{"initializing", "started"}, Dst: "failed"},
			},
			fsm.Callbacks{},
		),
	}
}

// Submission 은 유저의 제출 기록입니다. 소스 코드는 별도로 S3에 저장됩니다.
type Submission struct {
	ID   int64  `gorm:"primary_key" json:"id"`
	UUID string `gorm:"type:varchar(36);not null;unique_index" json:"uuid"`

	UserID int64  `gorm:"not null;index" json:"userID"`
	GameID string `gorm:"type:varchar(10);not null;index" json:"gameID"`

	// python3.6
	Runtime     string `gorm:"type:varchar(32);not null" json:"runtime"`
	Description string `gorm:"type:varchar(40)" json:"description"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewSubmission 은 새로운 제출을 만듭니다.
func NewSubmission(userID int64, gameID, runtime, desc string) *Submission {
	u, _ := uuid.NewRandom()
	return &Submission{
		UUID:        u.String(),
		GameID:      gameID,
		Runtime:     runtime,
		Description: desc,
	}
}

// GetS3Key 는 해당 제출 파일이 저장된 S3 키를 반환합니다.
func (s *Submission) GetS3Key() string {
	return "source/" + s.GameID + "/" + s.UUID
}
