# Battle Tetris

[![Go Version](https://img.shields.io/badge/Go-1.23.0-00ADD8?logo=go)](https://go.dev/)
[![WebSocket](https://img.shields.io/badge/WebSocket-Gorilla-blue)](https://github.com/gorilla/websocket)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

실시간 멀티플레이어 테트리스 게임 서버. Go 언어와 WebSocket을 활용한 고성능 게임 엔진으로 구현되었습니다.

## 목차

- [주요 기능](#주요-기능)
- [아키텍처](#아키텍처)
- [빠른 시작](#빠른-시작)
- [설정](#설정)
- [API 문서](#api-문서)
- [게임 규칙](#게임-규칙)
- [기술 스택](#기술-스택)
- [프로젝트 구조](#프로젝트-구조)

## 주요 기능

### 🎮 게임 플레이
- **실시간 멀티플레이어**: WebSocket 기반 실시간 대전
- **공격 시스템**: 2줄 이상 제거 시 상대방에게 공격 블록 전송
- **블록 미리보기**: 다음 블록 10개 미리 확인
- **난이도 증가**: 1분마다 블록 낙하 속도 100ms씩 증가 (최소 100ms)
- **점수 시스템**:
  - 라인 제거: 라인당 10점
  - 승리 보너스: 상대 플레이어 수 × 100점

### 🏆 경쟁 요소
- **랭킹 시스템**: 상위 20명 점수 기록 및 관리
- **대기실**: 플레이어 매칭 및 방 생성
- **관전 모드**: 진행 중인 게임 실시간 관전
- **닉네임 변경**: 게임 중 자유로운 닉네임 변경

### 🤖 AI 봇 지원
- **Beginner Bot**: 기본 플레이 봇
- **커스터마이징 가능**: BotAdapter를 통한 확장 가능

### 🔒 보안 및 안정성
- **동시성 제어**: sync.RWMutex를 활용한 데이터 레이스 방지
- **CORS 설정**: 개발/프로덕션 환경별 출처 제어
- **메시지 검증**: 발신자 검증 및 입력 유효성 검사
- **에러 처리**: 체계적인 로깅 및 에러 복구

## 아키텍처

### 시스템 구성도

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│   Browser   │ ◄─────► │  WebSocket Server │ ◄─────► │   Manager   │
│  (HTML/JS)  │         │   (Gin + Gorilla) │         │  (Game Pool)│
└─────────────┘         └──────────────────┘         └─────────────┘
                                  │                          │
                                  ▼                          ▼
                        ┌──────────────────┐       ┌─────────────┐
                        │  HandleMessage   │       │    Game     │
                        │   (Router)       │       │  (Engine)   │
                        └──────────────────┘       └─────────────┘
                                  │                          │
                        ┌─────────▼──────────┐              │
                        │   Bot Adapter      │ ◄────────────┘
                        │  (AI Interface)    │
                        └────────────────────┘
```

### 데이터 흐름

```
사용자 액션:
  User → WebSocket → Client.Read() → HandleMessage → Manager → Game.Ch → Game.run()

게임 이벤트:
  Game.run() → Manager.Ch → Manager.HandleMessage() → WebSocket.broadcast → Client.Write() → User

봇 액션:
  Bot → BotAdapter → WebSocket → HandleMessage → Manager → Game
```

### 핵심 컴포넌트

#### 1. WebSocket Server
- **역할**: 클라이언트 연결 관리 및 메시지 라우팅
- **주요 기능**:
  - 방(Room) 생성/삭제/관리
  - 클라이언트 등록/해제
  - 브로드캐스트 메시지 처리
- **동시성**: RWMutex로 rooms 맵 보호

#### 2. Manager
- **역할**: 게임 인스턴스 및 플레이어 관리
- **주요 기능**:
  - 게임 생성/시작/종료
  - 플레이어 등록/매칭
  - 게임 이벤트 처리 (승리, 패배, 점수)
  - 랭킹 관리
- **동시성**: RWMutex로 players 맵 보호

#### 3. Game
- **역할**: 테트리스 게임 로직 엔진
- **주요 기능**:
  - 블록 이동/회전/낙하
  - 라인 제거 및 점수 계산
  - 충돌 감지
  - 게임 오버 판정
- **동시성**:
  - RWMutex로 Cell 보드 보호
  - 고루틴 2개 운영 (run, autoDown)
  - 채널 기반 메시지 처리

#### 4. Client
- **역할**: WebSocket 클라이언트 래퍼
- **주요 기능**:
  - 메시지 송수신
  - 발신자 검증
  - 연결 관리
- **동시성**:
  - Read/Write 고루틴 분리
  - 채널 기반 안전한 종료

## 빠른 시작

### 요구 사항

- Go 1.23.0 이상
- 최신 웹 브라우저 (Chrome, Firefox, Safari, Edge)

### 설치

```bash
# 저장소 클론
git clone https://github.com/lmk/tetris.git
cd tetris

# 의존성 설치
go mod download

# 빌드
go build -o tetris
```

### 실행

```bash
# 기본 설정으로 실행 (포트 8090)
./tetris

# 커스텀 설정 파일 사용
./tetris -config=config.yaml

# 커맨드 라인 옵션
./tetris -port=8080 -log-info=true -log-debug=false
```

### 접속

브라우저에서 `http://localhost:8090` 접속

## 설정

### config.yaml 생성

```bash
cp config.yaml.example config.yaml
```

### 설정 파일 구조

```yaml
# 서버 설정
domain: localhost
port: 8090
https: false

# CORS 설정 (중요!)
allowed_origins:
  - "*"  # 개발 환경: 모든 출처 허용

# 로그 설정
log:
  datetime: false
  srcfile: true
  info: true
  warning: true
  error: true
  trace: true
  debug: false
```

### CORS 설정

#### 개발 환경
```yaml
allowed_origins:
  - "*"  # 모든 출처 허용
```

#### 프로덕션 환경
```yaml
allowed_origins:
  - "http://yourdomain.com"
  - "https://yourdomain.com"
  - "http://192.168.1.100"       # LAN IP
  - "http://192.168.1.100:8090"  # 포트 명시
```

### 다른 기기에서 접속

1. 서버 IP 확인
   ```bash
   # Linux/Mac
   ifconfig | grep "inet "

   # Windows
   ipconfig | findstr IPv4
   ```

2. `config.yaml`에 IP 추가
   ```yaml
   allowed_origins:
     - "http://192.168.1.100"
     - "http://192.168.1.100:8090"
   ```

3. 모바일/태블릿에서 접속
   ```
   http://192.168.1.100:8090
   ```

### 커맨드 라인 옵션

```bash
./tetris -h

옵션:
  -config string
        설정 파일 경로 (기본값: "config.yaml")
  -port int
        서버 포트 (기본값: 8090)
  -https
        HTTPS 모드 활성화
  -log-datetime
        로그에 날짜/시간 표시
  -log-srcfile
        로그에 소스 파일 위치 표시
  -log-info
        Info 레벨 로그 활성화
  -log-warning
        Warning 레벨 로그 활성화
  -log-error
        Error 레벨 로그 활성화
  -log-trace
        Trace 레벨 로그 활성화
  -log-debug
        Debug 레벨 로그 활성화
```

## API 문서

### WebSocket 엔드포인트

```
ws://localhost:8090/ws
```

### 메시지 프로토콜

모든 메시지는 JSON 형식으로 전송됩니다.

```go
type Message struct {
    Action       string      // 액션 타입
    Sender       string      // 발신자 닉네임
    RoomId       int         // 방 ID
    Data         string      // 추가 데이터
    Cells        [][]int     // 게임 보드 상태
    CurrentBlock *Block      // 현재 블록
    BlockIndexs  []int       // 블록 인덱스 배열
    Score        int         // 점수
    RoomList     []RoomInfo  // 방 목록
}
```

### 클라이언트 → 서버

| 액션 | 설명 | 파라미터 |
|------|------|----------|
| `set-nick` | 닉네임 변경 | `Data`: 새 닉네임 |
| `create-room` | 방 생성 | - |
| `join-room` | 방 참가 | `RoomId`: 방 ID |
| `leave-room` | 방 나가기 | - |
| `start-game` | 게임 시작 (방장만) | - |
| `block-rotate` | 블록 회전 | - |
| `block-left` | 블록 왼쪽 이동 | - |
| `block-right` | 블록 오른쪽 이동 | - |
| `block-down` | 블록 한 칸 내리기 | - |
| `block-drop` | 블록 즉시 낙하 | - |
| `get-rank` | 랭킹 조회 | `Data`: 조회할 개수 |
| `add-bot` | 봇 추가 (방장만) | - |

### 서버 → 클라이언트

| 액션 | 설명 | 데이터 |
|------|------|--------|
| `new-nick` | 닉네임 할당 | `Data`: 닉네임 |
| `set-nick` | 닉네임 변경 완료 | `Sender`: 변경한 사용자 |
| `create-room` | 방 생성 완료 | `RoomId`, `RoomList` |
| `join-room` | 방 참가 완료 | `RoomId`, `RoomList` |
| `leave-room` | 방 나감 | `Sender`: 나간 사용자 |
| `start-game` | 게임 시작 | `Cells`, `CurrentBlock`, `BlockIndexs` |
| `sync-game` | 게임 상태 동기화 | `Cells`, `CurrentBlock`, `Score` |
| `over-game` | 게임 오버 | `Sender`: 패배한 플레이어 |
| `end-game` | 게임 종료 (승리) | `Sender`: 승자, `Data`: "winner" 또는 "winner:순위", `Score` |
| `gift-full-blocks` | 공격 블록 전송 | `Sender`: 공격자, `Cells`: 공격 블록 |
| `erase-blocks` | 라인 제거 | `BlockIndexs`: 제거된 라인 번호 |
| `rank` | 랭킹 정보 | `RankList` |
| `error` | 에러 메시지 | `Data`: 에러 내용 |

## 게임 규칙

### 기본 규칙

1. **보드 크기**: 15행 × 10열
2. **블록 종류**: 7가지 테트로미노 (I, T, O, Z, L, J, S)
3. **점수 계산**:
   - 1줄 제거: 10점
   - 2줄 제거: 20점 (상대방에게 공격)
   - 3줄 제거: 30점 (상대방에게 공격)
   - 4줄 제거: 40점 (상대방에게 공격)

### 멀티플레이어 규칙

1. **공격 시스템**:
   - 2줄 이상 동시 제거 시 상대방에게 해당 라인 전송
   - 받은 블록은 하단에서 밀려 올라옴
   - 밀려난 상단 블록이 있으면 게임 오버

2. **승리 조건**:
   - 다른 모든 플레이어가 게임 오버될 때까지 생존
   - 승리 시 보너스: 100점 × (플레이어 수 - 1)

3. **게임 오버 조건**:
   - 새 블록이 생성될 공간이 없을 때
   - 공격 블록을 받아서 상단 블록이 밀려날 때

### 난이도 시스템

- **초기 속도**: 1000ms (1초)
- **속도 증가**: 1분마다 100ms씩 감소
- **최소 속도**: 100ms (0.1초)

## 기술 스택

### 백엔드
- **언어**: Go 1.23.0
- **웹 프레임워크**: Gin 1.9.1
- **WebSocket**: Gorilla WebSocket 1.5.0
- **설정 관리**: gopkg.in/yaml.v2
- **동시성**: sync.RWMutex, 채널 기반 통신

### 프론트엔드
- **HTML5/CSS3**: 반응형 UI
- **JavaScript**: 게임 렌더링 및 WebSocket 통신
- **Bootstrap**: UI 컴포넌트
- **jQuery**: DOM 조작

### 아키텍처 패턴
- **Pub/Sub**: 채널 기반 이벤트 브로드캐스팅
- **Actor Model**: 각 게임이 독립적인 고루틴으로 실행
- **Command Pattern**: 메시지 기반 액션 처리

## 프로젝트 구조

```
tetris/
├── main.go              # 애플리케이션 진입점
├── Config.go            # 설정 관리
├── Logger.go            # 로깅 시스템
├── Message.go           # 메시지 구조체 정의
│
├── WebsocketServer.go   # WebSocket 서버 (연결 관리)
├── HandleMessage.go     # 메시지 라우팅 및 처리
├── Client.go            # WebSocket 클라이언트 래퍼
│
├── Manager.go           # 게임 매니저 (플레이어/게임 관리)
├── Game.go              # 게임 엔진 (테트리스 로직)
├── Block.go             # 블록 구조체 및 로직
├── Cells.go             # 보드 셀 관리
│
├── RoomInfo.go          # 방 정보 구조체
├── Rank.go              # 랭킹 시스템
│
├── Bot.go               # 봇 인터페이스
├── BotAdapter.go        # 봇 어댑터 (WebSocket 연결)
├── BotBeginer.go        # 초급 봇 구현
├── BotFather.go         # 봇 관리자
│
├── config.yaml.example  # 설정 파일 예제
├── history.md           # 변경 이력
│
├── templates/           # HTML 템플릿
│   └── index.html
│
└── public/              # 정적 파일
    ├── tetris.js        # 게임 렌더링 로직
    ├── ws.js            # WebSocket 통신
    ├── tetris.css       # 게임 스타일
    ├── room.css         # 방 UI 스타일
    └── sound/           # 효과음
```

### 주요 파일 설명

#### 서버 코어
- **main.go** (106줄): 서버 초기화 및 라우팅 설정
- **WebsocketServer.go** (183줄): WebSocket 연결 및 방 관리
- **HandleMessage.go** (315줄): 클라이언트 요청 처리

#### 게임 로직
- **Game.go** (600줄): 테트리스 게임 엔진 핵심
- **Manager.go** (311줄): 게임 인스턴스 및 이벤트 관리
- **Block.go** (100줄): 테트로미노 정의 및 회전 로직
- **Cells.go** (133줄): 보드 셀 조작 유틸리티

#### 시스템
- **Config.go** (102줄): YAML 기반 설정 관리
- **Logger.go** (70줄): 레벨별 로깅 시스템
- **Rank.go** (136줄): 파일 기반 랭킹 저장

#### 봇 시스템
- **BotAdapter.go** (115줄): 봇과 서버 간 어댑터
- **BotBeginer.go** (130줄): AI 봇 구현 예제

### 코드 통계
- **총 라인 수**: ~3,144줄 (Go)
- **주요 고루틴**:
  - WebSocket 서버 실행 루프
  - Manager 이벤트 처리 루프
  - 게임당 2개 (run, autoDown)
  - 클라이언트당 2개 (Read, Write)
  - 봇당 2개 (Read, Write)

## 성능 최적화

### 동시성 제어
- **RWMutex**: 읽기 작업이 많은 맵에 사용하여 성능 향상
- **채널 버퍼링**: `MAX_CHAN` 크기 버퍼로 고루틴 블로킹 최소화
- **strings.Builder**: 문자열 연결 성능 최적화

### 메모리 관리
- **블록 재사용**: CloneShape로 불필요한 할당 방지
- **슬라이스 용량 예약**: make([]Type, 0, capacity) 패턴 사용
- **지연 초기화**: 필요한 시점에만 리소스 할당

### 네트워크 최적화
- **메시지 필터링**: 불필요한 브로드캐스트 최소화
- **압축 가능**: WebSocket 압축 확장 지원 가능
- **선택적 동기화**: 자신의 게임 상태만 전체 동기화

## 개발 가이드

### 새로운 봇 추가

1. `Bot` 인터페이스 구현
```go
type MyBot struct {
    Nick string
}

func (b *MyBot) GetNick() string {
    return b.Nick
}

func (b *MyBot) Think(cells [][]int, block *Block, nextBlocks []int) string {
    // AI 로직 구현
    return "block-drop"
}
```

2. `BotFather`에 등록
```go
func CreateBot(botType string, nick string) Bot {
    switch botType {
    case "mybot":
        return &MyBot{Nick: nick}
    // ...
    }
}
```

### 새로운 액션 추가

1. `Message` 구조체에 필드 추가 (필요시)
2. `HandleMessage.go`에 핸들러 추가
3. `Game.go`에 로직 추가
4. 프론트엔드 `ws.js`에 처리 로직 추가

### 테스트

```bash
# 유닛 테스트
go test ./...

# 특정 파일 테스트
go test -v Cells_test.go Cells.go

# 벤치마크
go test -bench=.
```

## 문제 해결

### 포트가 이미 사용 중
```bash
# Linux/Mac
lsof -ti:8090 | xargs kill -9

# Windows
netstat -ano | findstr :8090
taskkill /PID <PID> /F
```

### CORS 에러
1. `config.yaml`에서 `allowed_origins` 확인
2. 브라우저 개발자 도구에서 실제 Origin 확인
3. 해당 Origin을 allowed_origins에 추가

### WebSocket 연결 실패
1. 방화벽 설정 확인
2. 서버 로그에서 에러 확인 (`-log-error=true`)
3. 브라우저 콘솔에서 에러 메시지 확인

### 게임 오버 후 재시작 안됨
- 방장만 게임을 시작할 수 있습니다
- 방을 나갔다가 다시 들어와 방장이 되세요

## 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### 코딩 컨벤션
- Go 표준 포맷터 사용: `go fmt`
- Linter 통과: `golangci-lint run`
- 의미있는 커밋 메시지 작성
- 테스트 코드 작성 권장

## 라이선스

이 프로젝트는 MIT 라이선스 하에 있습니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 저자

**lmk** - [GitHub](https://github.com/lmk)

## 감사의 말

- Gin Web Framework
- Gorilla WebSocket
- Bootstrap
- jQuery

## 변경 이력

자세한 변경 이력은 [history.md](history.md)를 참조하세요.

### 최신 업데이트 (2025-12-22)

#### 버그 수정
- ✅ 승리 시 "Congratulations!!" 팝업이 제대로 표시되도록 수정
- ✅ 라인 제거 시 이상한 블록 생성 버그 수정
- ✅ 배열 범위 초과 버그 수정
- ✅ 동시성 데이터 레이스 수정
- ✅ 파일 쓰기 버그 수정
- ✅ 무한 루프 버그 수정

#### 개선 사항
- ✅ CORS 보안 개선 (모바일 접속 지원)
- ✅ 입력 검증 강화
- ✅ 성능 최적화 (strings.Builder)
- ✅ 코드 품질 개선 (상수화, 에러 처리)
- ✅ 문서화 개선

---

**즐거운 테트리스 대전 되세요!** 🎮✨
