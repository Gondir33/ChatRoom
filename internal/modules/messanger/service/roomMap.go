package service

/*
// Storage of Sockets rooms
type RoomMap struct {
	mutex sync.RWMutex
	rooms map[*websocket.Conn]int
}

func NewRoomMap() *RoomMap {
	return &RoomMap{
		rooms: make(map[*websocket.Conn]int, CapOfWebSockets),
	}
}

func (m *RoomMap) delete(conn *websocket.Conn) {
	m.mutex.Lock()
	delete(m.rooms, conn)
	m.mutex.Unlock()
}

func (m *RoomMap) set(conn *websocket.Conn, room int) {
	m.mutex.Lock()
	m.rooms[conn] = room
	m.mutex.Unlock()
}

func (m *RoomMap) findAllConnOfSuchRoom(room int) []*websocket.Conn {
	sockets := make([]*websocket.Conn, 0, CapOfSocketsRoom)
	m.mutex.RLock()
	for conn, r := range m.rooms {
		if r == room {
			sockets = append(sockets, conn)
		}
	}
	m.mutex.RUnlock()
	return sockets
}
*/
