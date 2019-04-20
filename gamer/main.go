package gamer

import (
	context "context"
	"errors"
	"strconv"
	"time"

	grpc "google.golang.org/grpc"
)

// ServerDriver 는 게이머 프로세스의 동작을 결정하는 인터페이스입니다.
// 각 언어별 구현체가 필요합니다.
type ServerDriver interface {
	StartProcess() error
}

// Gamer 는
type Gamer struct {
	UUID string

	client GamerClient
	// GetUUID() string

	// StartServer(port int) error
	// TakeAction(input []string) []string
}

// NewGamer 함수는 새로운 Gamer를 만듭니다.
func NewGamer(uuid string) *Gamer {
	return &Gamer{UUID: uuid}
}

// StartConnection 은 게이머 프로세스와의 커넥션을 맺습니다.
func (gamer *Gamer) StartConnection(port int, driver ServerDriver) error {
	driver.StartProcess()

	conn, err := grpc.Dial("127.0.0.1:"+strconv.Itoa(port), grpc.WithInsecure())

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)

		client := NewGamerClient(conn)
		_, err = client.Ping(context.Background(), &PingMessage{Id: "my_id"})
		if err != nil {
			if i == 9 {
				return errors.New("connection failed")
			}
			continue
		}

		gamer.client = client
		break
	}

	return nil
}

// TakeAction 함수는 Gamer에게 게임의 현 상태에 대한 input을 제공하고,
// 판단한 결과를 취합니다.
func (gamer *Gamer) TakeAction(input []string) ([]string, error) {
	actionID := "test"
	res, err := gamer.client.Action(context.Background(), &ActionInput{
		Id:   actionID,
		Data: input,
	})
	if err != nil {
		return nil, err
	}

	if actionID != res.Id {
		return nil, errors.New("invalid action id: " + res.Id)
	}
	return res.Data, nil
}

// Config 는
type Config struct {
	UUID string

	ServerPort int
	client     GamerClient
}
