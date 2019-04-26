# API Specification

## 인증

HTTP 요청 시 `Authorization` 헤더에 JWT 토큰을 넣어 요청합니다.

Token 생성 시 사용하는 payload는 다음과 같습니다.

```json
{
  "access_token": "YOUR_ACCESS_TOKEN"
}
```
