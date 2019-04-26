package game

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/hellodhlyn/undersky/gamer"
)

// Game1000 - 틱택토
//
// 입력 문자열
//   [0] 선공/후공 여부 (1 = 선공, 2 = 후공)
//   [1] 현재의 게임판 상황을 표시한 3행 3열의 문자열
//       0 = 아무것도 표시되어있지 않음, 1 = 선공, 2 = 후공
//       ex> 000
//           010
//           002
//
// 출력 문자열
//   [0] 내 차례에 표시할 좌표
//       알파벳 한 글자, 숫자 한 글자의 두 자리 문자열 (왼쪽부터 A, B, C, 위에서부터 1, 2, 3)
type Game1000 struct {
	ctx   *MatchContext
	board *TicTacToeBoard

	firstGamer  *gamer.Gamer
	secondGamer *gamer.Gamer
}

// GetRuleset 함수는 게임의 룰을 반환합니다.
//   - 9판 5선승제
func (*Game1000) GetRuleset() *Ruleset {
	return &Ruleset{
		MaximumRound: 9,
		RoundToWin:   5,
	}
}

// InitMatch 함수는 게임의 기본을 설정합니다.
func (game *Game1000) InitMatch(ctx *MatchContext) error {
	game.ctx = ctx

	return nil
}

// InitRound 함수는 각 라운드별 데이터를 초기화합니다.
func (game *Game1000) InitRound() error {
	game.board = NewTicTacToeBoard()

	// 선공과 후공을 정합니다. 선공은 1, 후공은 2를 놓습니다.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(2) == 0 {
		game.firstGamer, game.secondGamer = game.ctx.Player, game.ctx.Competitor
	} else {
		game.firstGamer, game.secondGamer = game.ctx.Competitor, game.ctx.Player
	}

	return nil
}

// PlayRound 함수는 라운드를 실행합니다.
func (game *Game1000) PlayRound() (string, error) {
	step := 1
	for {
		// 차례를 결정합니다.
		var actor *gamer.Gamer
		if step%2 == 1 {
			actor = game.firstGamer
		} else {
			actor = game.secondGamer
		}

		var ruleBreakError error
		if actor.UUID == game.ctx.Player.UUID {
			ruleBreakError = ErrPlayerBreakRule
		} else {
			ruleBreakError = ErrCompetitorBreakRule
		}

		// 동작을 취합니다.
		output, err := actor.TakeAction([]string{strconv.Itoa((step+1)%2 + 1), game.board.GetInputText()})
		if err != nil {
			return "", ErrInternalError
		}
		if len(output) != 1 || len(output[0]) != 2 {
			return "", ruleBreakError
		}
		game.board.Set(output[0], int8((step+1)%2+1))

		switch game.board.FindWinner() {
		case 1:
			return game.firstGamer.UUID, nil
		case 2:
			return game.secondGamer.UUID, nil
		}

		if step == 9 {
			break
		}
		step++
	}

	return "", nil
}

// TicTacToeBoard 는 틱택토 게임을 진행하는 보드입니다.
type TicTacToeBoard struct {
	board [][]int8
}

// NewTicTacToeBoard 함수는 새로운 보드를 만듭니다.
func NewTicTacToeBoard() *TicTacToeBoard {
	return &TicTacToeBoard{
		board: [][]int8{
			[]int8{0, 0, 0},
			[]int8{0, 0, 0},
			[]int8{0, 0, 0},
		},
	}
}

// GetInputText 함수는 현재 플레이 상황을 표시하는 문자열을 반환합니다.
func (board TicTacToeBoard) GetInputText() string {
	boardText := ""
	for idx, row := range board.board {
		for _, col := range row {
			boardText = boardText + strconv.Itoa(int(col))
		}

		if idx != 2 {
			boardText = boardText + "\n"
		}
	}
	return boardText
}

// Set 함수는 보드에 특정 값을 표기합니다.
func (board TicTacToeBoard) Set(position string, value int8) error {
	yAxis, err := strconv.Atoi(string(position[1]))
	if err != nil || yAxis < 1 || yAxis > 3 {
		return errors.New("invalid position")
	}
	yAxis = yAxis - 1

	xAxis := 0
	switch position[0] {
	case 'A':
		xAxis = 0
	case 'B':
		xAxis = 1
	case 'C':
		xAxis = 2
	default:
		return errors.New("invalid position")
	}

	if board.board[yAxis][xAxis] != 0 {
		return errors.New("invalid position")
	}

	board.board[yAxis][xAxis] = value

	return nil
}

// FindWinner 함수는 승리 여부를 검사합니다.
func (board TicTacToeBoard) FindWinner() int8 {
	for i := 0; i < 3; i++ {
		if board.board[0][i] == board.board[1][i] && board.board[1][i] == board.board[2][i] && board.board[0][i] != 0 {
			return board.board[0][i]
		} else if board.board[i][0] == board.board[i][1] && board.board[i][1] == board.board[i][2] && board.board[i][0] != 0 {
			return board.board[i][0]
		}
	}

	if board.board[0][0] == board.board[1][1] && board.board[1][1] == board.board[2][2] && board.board[0][0] != 0 {
		return board.board[0][0]
	} else if board.board[0][2] == board.board[1][1] && board.board[1][1] == board.board[2][0] && board.board[0][2] != 0 {
		return board.board[0][2]
	}
	return 0
}
