package game

import (
	"errors"

	"github.com/hellodhlyn/undersky/gamer"
)

var (
	// ErrInvalidGameContext 오류는 게임의 정보가 잘못되었을 때 반환됩니다.
	ErrInvalidGameContext = errors.New("invalid game context")

	// ErrInternalError 오류는 게임 실행 중 내부적인 오류가 발생했을 떄 반환됩니다.
	ErrInternalError = errors.New("unexpected error occurred")

	// ErrPlayerBreakRule 오류는 유저(제출자)가 잘못된 동작을 취했을 때 반환됩니다.
	ErrPlayerBreakRule = errors.New("player take invalid action")

	// ErrCompetitorBreakRule 오류는 상대가 잘못된 동작을 취했을 때 반환됩니다.
	ErrCompetitorBreakRule = errors.New("competitor take invalid action")
)

// Ruleset 은 게임에 관한 규칙에 대한 타입입니다.
type Ruleset struct {
	MaximumRound int
	RoundToWin   int
}

// MatchContext 는 현재 진행중인 경기에 대한 컨텍스트입니다.
type MatchContext struct {
	MatchUUID string

	Player     *gamer.Gamer
	Competitor *gamer.Gamer
}

// RoundResult 는 플레이한 라운드의 결과 타입입니다.
type RoundResult struct {
	WinnerID int64
}

// Game 은 게임 플레이에 대한 인터페이스입니다.
//
// 호출 순서
//   - InitMatch()
//   - InitRound()
//   - PlayRound()
type Game interface {
	GetRuleset() *Ruleset

	InitMatch(*MatchContext) error
	InitRound() error
	PlayRound() (*RoundResult, error)
}
