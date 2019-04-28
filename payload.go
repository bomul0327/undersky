package undersky

import "encoding/json"

// SubmissionPayload 는 Submission을 SQS로 주고받을 때 사용하는 페이로드입니다.
type SubmissionPayload struct {
	MatchUUID string `json:"matchUUID"`
}

// ToJSON 은 해당 인스턴스를 JSON byte 문자열로 변환합니다.
func (p *SubmissionPayload) ToJSON() []byte {
	j, _ := json.Marshal(p)
	return j
}
