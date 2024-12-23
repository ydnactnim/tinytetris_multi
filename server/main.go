package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket 업그레이더
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 플레이어 구조체
type Player struct {
	ID      string
	Name    string
	Conn    *websocket.Conn
	IsReady bool
	Score   int
	Field   [20][10]int
}

// 게임 상태 구조체
type GameState struct {
	PlayerID     string      `json:"playerID"`
	Field        [20][10]int `json:"field"`
	CurrentBlock string      `json:"currentBlock"`
	Score        int         `json:"score"`
}

// 방 구조체
type Room struct {
	ID          string
	Players     map[string]*Player
	Mutex       sync.Mutex
	GameStarted bool
}

// 서버 구조체
type Server struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

// 새 서버 생성
func NewServer() *Server {
	return &Server{
		Rooms: make(map[string]*Room),
	}
}

// 방 생성
func (s *Server) CreateRoom(roomID string) *Room {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room := &Room{
		ID:      roomID,
		Players: make(map[string]*Player),
	}
	s.Rooms[roomID] = room
	return room
}

// 방 참가
func (s *Server) JoinRoom(roomID, playerID string, player *Player) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 방이 없으면 자동 생성
	room, exists := s.Rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			Players: make(map[string]*Player),
		}
		s.Rooms[roomID] = room
		fmt.Println("Room created:", roomID)
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	room.Players[playerID] = player
	fmt.Printf("Player %s joined room %s\n", playerID, roomID)
	return nil
}

// WebSocket 연결 처리
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	// WebSocket 업그레이드
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// 플레이어 등록
	playerID := r.URL.Query().Get("playerID")
	playerName := r.URL.Query().Get("name")
	player := &Player{
		ID:   playerID,
		Name: playerName,
		Conn: conn,
	}

	// 방 참가
	roomID := r.URL.Query().Get("roomID")
	err = s.JoinRoom(roomID, playerID, player)
	if err != nil {
		fmt.Println("Join room error:", err)
		return
	}

	// 연결 종료 시 플레이어 제거
	defer func() {
		room, exists := s.Rooms[roomID]
		if exists {
			room.Mutex.Lock()
			delete(room.Players, playerID)
			room.Mutex.Unlock()
			fmt.Printf("Player %s left room %s\n", playerID, roomID)
		}
	}()

	// 메시지 수신 처리
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		switch msg["type"] {
		case "update":
			// 플레이어의 게임 상태 업데이트
			state := GameState{
				PlayerID:     player.ID,
				Field:        msg["field"].([20][10]int),
				CurrentBlock: msg["currentBlock"].(string),
				Score:        int(msg["score"].(float64)),
			}
			s.BroadcastGameState(roomID, state)
		default:
			fmt.Println("Unknown message type:", msg["type"])
		}
	}
}

// Ready 상태 변경
func (s *Server) SetPlayerReady(roomID, playerID string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room, exists := s.Rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		return fmt.Errorf("player not found in room")
	}

	player.IsReady = true
	fmt.Printf("Player %s in room %s is ready\n", playerID, roomID)

	// 모든 플레이어가 Ready 상태인지 확인
	allReady := true
	for _, p := range room.Players {
		if !p.IsReady {
			allReady = false
			break
		}
	}

	if allReady && !room.GameStarted {
		room.GameStarted = true
		fmt.Printf("All players in room %s are ready. Starting game!\n", roomID)

		// 초기 상태 브로드캐스트
		for _, p := range room.Players {
			initialState := GameState{
				PlayerID:     p.ID,
				Field:        [20][10]int{}, // 초기 빈 필드
				CurrentBlock: "I",           // 임의로 초기 블록 설정
				Score:        0,
			}
			s.BroadcastGameState(roomID, initialState)
		}
	}

	return nil
}

// 게임 상태 브로드캐스트
func (s *Server) BroadcastGameState(roomID string, state GameState) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room, exists := s.Rooms[roomID]
	if !exists {
		fmt.Println("Room not found for broadcasting game state")
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	// 각 플레이어에게 상태 전송
	for _, player := range room.Players {
		if player.ID != state.PlayerID { // 상태 전송 제외: 자기 자신
			err := player.Conn.WriteJSON(state)
			if err != nil {
				fmt.Printf("Error broadcasting to player %s: %v\n", player.ID, err)
			}
		}
	}
}

// 주기적 상태 업데이트
func (s *Server) PeriodicStateUpdate() {
	for {
		time.Sleep(100 * time.Millisecond)
		s.Mutex.Lock()
		for _, room := range s.Rooms {
			room.Mutex.Lock()
			if room.GameStarted {
				// 예시: 각 플레이어의 상태를 주기적으로 브로드캐스트
				for _, player := range room.Players {
					state := GameState{
						PlayerID:     player.ID,
						Field:        player.Field,
						CurrentBlock: "exampleBlock", // 임시 블록 값
						Score:        player.Score,
					}
					s.BroadcastGameState(room.ID, state)
				}
			}
			room.Mutex.Unlock()
		}
		s.Mutex.Unlock()
	}
}

// 서버 시작
func main() {
	server := NewServer()

	http.HandleFunc("/ws", server.handleConnection)

	// 주기적 상태 업데이트
	go server.PeriodicStateUpdate()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
