# Battle Tetris

웹 기반 대전 테트리스 게임 서버 (Go + WebSocket)

## Features

- Waiting room
- Preview
- Multi play
- Change nick
- Rank

## Configuration

### CORS 설정 (중요!)

서버를 다른 기기에서 접속하려면 `config.yaml` 파일에서 CORS 설정이 필요합니다.

#### 개발 환경 (모든 출처 허용)
```yaml
allowed_origins:
  - "*"
```

#### 프로덕션 환경 (특정 출처만 허용)
```yaml
allowed_origins:
  - "http://yourdomain.com"
  - "https://yourdomain.com"
  - "http://192.168.1.100"      # 서버 IP
  - "http://192.168.1.100:8090" # 포트 포함
```

### 설정 파일 생성

1. `config.yaml.example`을 복사하여 `config.yaml` 생성
2. 필요에 따라 설정 수정

```bash
cp config.yaml.example config.yaml
```

## Quick Start

```bash
# 빌드
go build

# 실행 (기본 포트: 8090)
./tetris

# 설정 파일 지정
./tetris -config=config.yaml
```

### 다른 기기에서 접속하기

1. 서버의 IP 주소 확인 (예: 192.168.1.100)
2. `config.yaml`에 해당 IP를 `allowed_origins`에 추가
3. 핸드폰/태블릿 브라우저에서 `http://192.168.1.100:8090` 접속

## TODO

- BOT
  - bigenner: 한블럭씩 & 손 느림
  - Light Finger: 한 불럭씩 손이 빠름
  - attacker2: 두블럭씩 쌓아서 공격
  - attacker3: 세블럭씩 쌓아서 공격
  - attacker half: 보드의 반을 쌓아서 공격

## flow

### call

user - html/js - WebsocketServer - HandleMessage - Manager - Game

bot - BotAdapter - WebsocketServer - HandleMessage - Manager - Game

### sequence
