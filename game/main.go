package game

import (
	"errors"

	"github.com/hellodhlyn/undersky-colosseum/gamer"
)

var (
	// ErrInvalidGameContext 오류는 게임의 정보가 잘못되었을 때 반환됩니다.
	ErrInvalidGameContext = errors.New("invalid game context")

	// ErrInternalError 오류는 게임 실행 중 내부적인 오류가 발생했을 떄 반환됩니다.
	ErrInternalError = errors.New("unexpected error occurred")

	// ErrPlayerBreakRule 오류는 유저(제출자)가 잘못된 동작을 취했을 때 반환됩니다.
	ErrPlayerBreakRule = errors.New("player take invalid action")

	// ErrCompetitionBreakRule 오류는 상대가 잘못된 동작을 취했을 때 반환됩니다.
	ErrCompetitionBreakRule = errors.New("competition take invalid action")
)

// Ruleset 은 게임에 관한 규칙에 대한 타입입니다.
type Ruleset struct {
	MaximumRound int
	RoundToWin   int
}

// Context 는
type Context struct {
	GameID int64

	Player      *gamer.Gamer
	Competition *gamer.Gamer
}

// Game 은 게임 플레이에 대한 인터페이스입니다.
//
// 호출 순서
//   - InitGame()
//   - InitRound()
//   - PlayRound()
type Game interface {
	GetRuleset() *Ruleset

	InitGame(*Context) error
	InitRound() error
	PlayRound() (string, error)
}
