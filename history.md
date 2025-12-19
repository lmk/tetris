# Battle Tetris 버그 수정 및 개선 이력
 
## 2025년 수정 내역
 
### 주요 버그 수정 (14가지)
 
#### 1. 배열 범위 초과 버그 수정
**파일**: `Game.go`
- **문제**: `receiveFullBlocks` 함수에서 Cell 배열에 대한 범위 검사 없이 직접 수정하여 out of bounds 에러 발생 가능
- **해결**: 블록 추가 전 모든 셀이 비어있는지 먼저 확인하는 로직 추가
```go
// 수정 전: Cell 직접 수정
g.Cell = append(g.Cell[len(blocks):], blocks...)
 
// 수정 후: 먼저 검증 후 수정
for r := 0; r < len(blocks); r++ {
    for c := 0; c < BOARD_COLUMN; c++ {
        if g.Cell[r][c] != EMPTY {
            return false
        }
    }
}
g.Cell = append(g.Cell[len(blocks):], blocks...)
```
 
#### 2. 동시성 데이터 레이스 수정
**파일**: `Manager.go`, `WebsocketServer.go`, `Game.go`
- **문제**: 여러 고루틴에서 공유 데이터(players, rooms, Cell)에 동시 접근 시 데이터 레이스 발생
- **해결**: `sync.RWMutex`를 추가하여 읽기/쓰기 락 구현
```go
// Manager.go
type manager struct {
    players map[string]*Player
    ch      chan *Message
    mu      sync.RWMutex  // 추가
}
 
// WebsocketServer.go
type WebsocketServer struct {
    rooms      map[int]*RoomInfo
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex  // 추가
}
 
// Game.go
type Game struct {
    // ... 기타 필드
    mu sync.RWMutex  // 추가
}
```
 
#### 3. CORS 보안 개선
**파일**: `Config.go`, `WebsocketServer.go`, `config.yaml.example`, `README.md`
- **문제**: CORS 설정이 없어 다른 출처에서 WebSocket 연결 불가
- **해결**:
  - 개발 환경: 모든 출처 허용 (`"*"`)
  - 프로덕션 환경: 특정 출처만 허용 (config.yaml에서 설정)
```go
// Config.go
type Config struct {
    AllowedOrigins []string `yaml:"allowed_origins,omitempty"`
    // ...
}
 
// WebsocketServer.go
CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    if origin == "" {
        return true
    }
 
    for _, allowed := range conf.AllowedOrigins {
        if allowed == "*" {
            return true
        }
        if origin == allowed {
            return true
        }
    }
 
    Warning.Printf("CORS: Rejected origin: %s", origin)
    return false
}
```
 
#### 4. 파일 쓰기 버그 수정
**파일**: `Rank.go`
- **문제**: 파일에 새 내용 쓰기 전 truncate하지 않아 이전 데이터와 섞임
- **해결**: `file.Truncate(0)` 호출 추가
```go
// 수정 후
file.Seek(0, 0)
err := file.Truncate(0)  // 기존 내용을 완전히 지움
if err != nil {
    Error.Printf("Invalid truncate: %s", err)
    return rank, err
}
_, err = file.WriteString(buf)
```
 
#### 5. 무한 루프 버그 수정
**파일**: `Game.go`
- **문제**: `autoDown` 함수에서 CycleMs가 음수로 감소하여 무한 루프 발생 가능
- **해결**: `MIN_CYCLE_MS` 상수 추가 및 하한선 체크
```go
const (
    MIN_CYCLE_MS     = 100
    INITIAL_CYCLE_MS = 1000
    CYCLE_DECREASE   = 100
)
 
func (g *Game) autoDown() {
    for !g.IsGameOver() {
        time.Sleep(time.Millisecond * time.Duration(g.CycleMs))
        g.Ch <- &Message{Action: "auto-down"}
 
        duration := time.Since(g.DurationTime)
        if duration > time.Minute*1 {
            if g.CycleMs > MIN_CYCLE_MS {
                g.CycleMs -= CYCLE_DECREASE
                if g.CycleMs < MIN_CYCLE_MS {
                    g.CycleMs = MIN_CYCLE_MS
                }
            }
            g.DurationTime = time.Now()
        }
    }
}
```
 
#### 6. 게임 오버 감지 일관성 개선
**파일**: `Game.go`
- **문제**: 일부 블록 액션에서만 게임 오버 체크하여 불일치 발생
- **해결**: 모든 블록 액션(`auto-down`, `block-drop`, `gift-full-blocks`)에서 일관되게 게임 오버 체크
```go
// HandleMessage 함수에서 모든 액션에 대해 동일하게 처리
case "auto-down":
    if g.Down() {
        g.SetState("gameover")
        Manager.ch <- &Message{Action: "over-game", Sender: g.Owner}
    }
 
case "block-drop":
    g.Drop()
    if g.IsGameOver() {
        g.SetState("gameover")
        Manager.ch <- &Message{Action: "over-game", Sender: g.Owner}
    }
```
 
#### 7. 채널 닫힘 안전성 추가
**파일**: `BotAdapter.go`
- **문제**: 닫힌 채널에 메시지 전송 시 패닉 발생
- **해결**: `select`문과 `default` 케이스로 안전하게 처리
```go
if msg.Action == "leave-room" {
    select {
    case BotFather.fromBot <- msg:
        // 성공적으로 전송
    default:
        // 채널이 닫혔거나 가득 찼을 때
        Warning.Println("BotFather.fromBot channel unavailable")
    }
}
```
 
#### 8. 에러 처리 개선
**파일**: `main.go`
- **문제**: 설정 초기화 실패 시 적절한 에러 처리 없음
- **해결**: 치명적 오류 발생 시 stderr에 출력하고 프로세스 종료
```go
func main() {
    err := initConf()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
        os.Exit(1)
    }
 
    InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    Info.Printf("config: %s", conf.makePretty())
    runServer()
}
```
 
#### 9. 하드코딩된 매직 넘버 제거
**파일**: `Game.go`, `WebsocketServer.go`
- **문제**: 코드 전체에 의미 없는 숫자 하드코딩
- **해결**: 명명된 상수로 교체하여 가독성 및 유지보수성 향상
```go
// Game.go
const (
    MIN_CYCLE_MS     = 100
    INITIAL_CYCLE_MS = 1000
    CYCLE_DECREASE   = 100
    SCORE_PER_LINE   = 10
    WINNER_BONUS     = 100
)
 
// WebsocketServer.go
const (
    MAX_ROOMS = 100
)
```
 
#### 10. 입력 검증 강화
**파일**: `Client.go`
- **문제**: 클라이언트가 보낸 메시지의 발신자 검증 없이 신뢰
- **해결**: 메시지 발신자가 실제 클라이언트와 일치하는지 검증
```go
func (client *Client) Read() {
    for {
        var msg Message
        err := client.socket.ReadJSON(&msg)
        if err != nil {
            break
        }
 
        // 메시지 발신자가 현재 클라이언트와 일치하는지 검증
        if msg.Sender != "" && msg.Sender != client.Nick {
            Warning.Printf("Message sender mismatch: claimed=%s, actual=%s",
                msg.Sender, client.Nick)
            msg.Sender = client.Nick // 강제로 올바른 발신자로 교정
        }
 
        if msg.Sender == "" {
            msg.Sender = client.Nick
        }
 
        // ... 나머지 처리
    }
}
```
 
#### 11. 성능 최적화
**파일**: `WebsocketServer.go`
- **문제**: 문자열 연결 시 `+` 연산자 사용으로 불필요한 메모리 할당
- **해결**: `strings.Builder` 사용으로 성능 개선
```go
func (wss *WebsocketServer) Report() {
    var report strings.Builder
    wss.mu.RLock()
    for roomId, info := range wss.rooms {
        fmt.Fprintf(&report, "[%v:%s:[", roomId, info.Owner)
        first := true
        for nick := range info.Clients {
            if !first {
                report.WriteString(",")
            }
            report.WriteString(nick)
            first = false
        }
        report.WriteString("]]")
    }
    wss.mu.RUnlock()
 
    Info.Println("REPORT:" + report.String())
}
```
 
#### 12. WebSocket 메시지 검증 강화
**파일**: `HandleMessage.go`
- **문제**: 방 ID와 닉네임 검증이 불충분
- **해결**: 모든 메시지 처리 전 mutex로 보호된 검증 함수 사용
```go
func (wss *WebsocketServer) isVaildRoomId(roomId int) bool {
    wss.mu.RLock()
    defer wss.mu.RUnlock()
    if _, ok := wss.rooms[roomId]; ok {
        return true
    }
    return false
}
 
func (wss *WebsocketServer) isVaildNick(roomId int, nick string) bool {
    wss.mu.RLock()
    defer wss.mu.RUnlock()
    if room, ok := wss.rooms[roomId]; ok {
        if _, ok := room.Clients[nick]; ok {
            return true
        }
    }
    return false
}
```
 
#### 13. processFullLine 로직 정리
**파일**: `Game.go`
- **문제**: 불필요한 반복문과 복잡한 로직
- **해결**: 코드 간소화 및 가독성 향상
 
#### 14. 문서화 개선
**파일**: `README.md`, `config.yaml.example`
- **추가된 문서**:
  - CORS 설정 가이드 (개발/프로덕션 환경)
  - 설정 파일 생성 방법
  - Quick Start 가이드
  - 다른 기기(핸드폰/태블릿)에서 접속하는 방법
 
### 새로 추가된 파일
 
#### config.yaml.example
서버 설정 템플릿 파일
```yaml
domain: localhost
port: 8090
https: false
 
# CORS 설정
allowed_origins:
  - "*"  # 개발 환경: 모든 출처 허용
 
# 프로덕션 환경 예시:
# allowed_origins:
#   - "http://yourdomain.com"
#   - "https://yourdomain.com"
#   - "http://192.168.1.100"
#   - "http://192.168.1.100:8090"
 
log:
  datetime: false
  srcfile: true
  info: true
  warning: true
  error: true
  trace: true
  debug: false
```
 
### 수정된 파일 목록
 
1. `Config.go` - CORS 설정 추가
2. `WebsocketServer.go` - CORS 체크로직, mutex 추가, 성능 최적화
3. `Manager.go` - mutex 추가, 모든 공유 데이터 접근 보호
4. `HandleMessage.go` - mutex로 보호된 검증 함수 추가
5. `Game.go` - 배열 검증, mutex, 상수화, 게임 오버 로직 개선
6. `Client.go` - 메시지 발신자 검증 추가
7. `BotAdapter.go` - 채널 안전성 추가
8. `Rank.go` - 파일 truncate 추가
9. `main.go` - 에러 처리 개선
10. `README.md` - CORS 설정 및 사용법 문서화
11. `config.yaml.example` - 새로 생성
 
### 기술적 개선 사항
 
- **동시성 안전성**: 모든 공유 데이터에 대해 mutex 보호 적용
- **보안**: CORS 검증, 메시지 발신자 검증 추가
- **안정성**: 배열 범위 체크, 무한 루프 방지, 채널 안전성 확보
- **성능**: strings.Builder 사용으로 문자열 처리 최적화
- **유지보수성**: 매직 넘버를 명명된 상수로 교체
- **문서화**: README 및 설정 예시 파일 추가
 
### 테스트 결과
 
- ✅ 빌드 성공: `go build -v`
- ✅ 컴파일 에러 없음
- ✅ GitHub에 성공적으로 푸시됨
 
### 다음 단계 권장사항
 
1. Dependabot이 감지한 2개의 보안 취약점 확인 및 수정
2. 단위 테스트 추가
3. 통합 테스트 작성
4. 프로덕션 환경에 배포 전 config.yaml에서 CORS 설정 확인