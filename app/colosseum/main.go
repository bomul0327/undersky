package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hellodhlyn/undersky/libs/s3"

	"github.com/hellodhlyn/sqstask"

	"github.com/hellodhlyn/undersky/game"
	"github.com/hellodhlyn/undersky/gamer"
)

var games = map[string]game.Game{
	"1000": &game.Game1000{},
}

type matchOptionsGamer struct {
	UUID       string `json:"uuid"`
	Runtime    string `json:"runtime"`
	SourceUUID string `json:"source_uuid"`
}

type matchOptions struct {
	MatchID    int64             `json:"match_id"`
	GameID     string            `json:"game_id"`
	Player     matchOptionsGamer `json:"player"`
	Competitor matchOptionsGamer `json:"competitor"`
}

func runMatch(opts *matchOptions) {
	g, ok := games[opts.GameID]
	if !ok {
		fmt.Printf("invalid game id: %s\n", opts.GameID)
		return
	}

	// 게이머들의 프로세스를 실행합니다.
	fmt.Println("waiting for player...")
	player, err := makeGamer(opts.GameID, opts.Player)
	if err != nil {
		fmt.Printf("failed to create player: %v\n", err)
		return
	}

	fmt.Println("waiting for competitor...")
	competitor, err := makeGamer(opts.GameID, opts.Competitor)
	if err != nil {
		fmt.Printf("failed to create competitor: %v\n", err)
		return
	}

	fmt.Println("initializing game...")
	matchCtx := game.MatchContext{
		GameID:     opts.MatchID,
		Player:     player,
		Competitor: competitor,
	}
	g.InitMatch(&matchCtx)

	// 게임을 시작합니다.
	var playerWins int8
	var competitorWins int8
	for i := 0; i < g.GetRuleset().MaximumRound; i++ {
		fmt.Println("initializing round...")
		g.InitRound()

		fmt.Println("starting round...")
		winner, err := g.PlayRound()
		if err != nil {
			fmt.Printf("error on playing round: %v\n", err)
			return
		}

		if winner == player.UUID {
			fmt.Println("player win")
			playerWins++
		} else if winner == competitor.UUID {
			fmt.Println("competitor win")
			competitorWins++
		} else {
			fmt.Println("draw")
		}
	}

	fmt.Printf("[Result] Player %d : %d Competitor\n", playerWins, competitorWins)
}

func makeGamer(gameID string, opts matchOptionsGamer) (*gamer.Gamer, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	port := 10000 + r.Int()%55535

	g := gamer.NewGamer(opts.UUID)

	var driver gamer.ServerDriver
	switch opts.Runtime {
	case "python3.6":
		s3client, err := s3.NewClient("undersky-ai")
		if err != nil {
			return nil, err
		}

		err = s3client.Download("source/"+gameID+"/"+opts.SourceUUID, "source/"+strconv.Itoa(port)+".py")
		if err != nil {
			return nil, err
		}

		driver = gamer.NewPython3Driver("source." + strconv.Itoa(port))

	default:
		return nil, errors.New("no such runtime: " + opts.Runtime)
	}

	return g, g.StartConnection(port, driver)
}

func main() {
	msg := flag.String("message", "", "sqs message for debug")
	flag.Parse()
	if *msg != "" {
		var opts matchOptions
		json.Unmarshal([]byte(*msg), &opts)
		runMatch(&opts)
		return
	}

	task, _ := sqstask.NewSQSTask(&sqstask.Options{
		QueueName:  "undersky-submission",
		AWSRegion:  "ap-northeast-2",
		WorkerSize: 1,
		Consume: func(message string) error {
			var opts matchOptions
			json.Unmarshal([]byte(message), &opts)
			runMatch(&opts)
			return nil
		},
		HandleError: func(err error) {
			fmt.Println(err)
		},
	})

	task.StartConsumer()
}
