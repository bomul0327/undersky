package undersky

import (
	"database/sql/driver"
	"errors"

	"github.com/looplab/fsm"
)

// FSM 은 상태 머신 자료형입니다.
// TODO - looplab/fsm 라이브러리를 fork하여 직접 Value, Scan 함수 구현
type FSM struct {
	State *fsm.FSM
}

func NewFSM(initial string, events []fsm.EventDesc, callbacks map[string]fsm.Callback) *FSM {
	return &FSM{
		fsm.NewFSM(initial, events, callbacks),
	}
}

func (m FSM) Value() (driver.Value, error) {
	if m.State == nil {
		return nil, nil
	}
	return m.State.Current(), nil
}

func (m *FSM) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch src.(type) {
	case string:
		m.State.SetState(src.(string))
	case []byte:
		m.State.SetState(string(src.([]byte)))
	default:
		return errors.New("incompatible type for FSM")
	}

	return nil
}

func (m *FSM) Event(event string, args ...interface{}) error {
	return m.State.Event(event, args)
}
