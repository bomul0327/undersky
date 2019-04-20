package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/hellodhlyn/undersky-colosseum/game"
	"github.com/hellodhlyn/undersky-colosseum/gamer"
)

var games = map[string]game.Game{
	"1000": &game.Game1000{},
}

func main() {
	// 게임을 설정합니다.
	exeID := flag.Int64("id", -1, "id of game execution")
	gameNum := flag.String("game", "", "number of game")

	flag.Parse()

	if *exeID == -1 {
		panic("invalid execution id: " + strconv.FormatInt(*exeID, 10))
	}

	g, ok := games[*gameNum]
	if !ok {
		panic("no such game: " + *gameNum)
	}

	// 게이머들의 프로세스를 실행합니다.
	player := gamer.NewGamer("00000000-0000-0000-0000-000000000000")
	fmt.Println("waiting for player...")
	if err := player.StartConnection(50051, &gamer.Python3Driver{}); err != nil {
		panic(err)
	}

	competition := gamer.NewGamer("00000000-0000-0000-0000-000000000001")
	fmt.Println("waiting for competition...")
	if err := competition.StartConnection(50052, &gamer.Python3Driver{}); err != nil {
		panic(err)
	}

	fmt.Println("initializing game...")
	gameCtx := game.Context{
		GameID:      *exeID,
		Player:      player,
		Competition: competition,
	}
	g.InitGame(&gameCtx)

	// 게임을 시작합니다.
	var playerWins int8
	var competitionWins int8
	for i := 0; i < g.GetRuleset().MaximumRound; i++ {
		fmt.Println("initializing round...")
		g.InitRound()

		fmt.Println("starting round...")
		winner, err := g.PlayRound()
		if err != nil {
			panic(err)
		}

		if winner == player.UUID {
			fmt.Println("player win")
			playerWins++
		} else if winner == competition.UUID {
			fmt.Println("competition win")
			competitionWins++
		} else {
			fmt.Println("draw")
		}
	}

	fmt.Printf("[Result] Player %d : %d Competition\n", playerWins, competitionWins)
}
