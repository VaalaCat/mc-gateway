package mc

const (
	ErrNotJoined = "not joined"
)

const (
	EventPlayerJoined = "multiplayer.player.joined"
	EventPlayerLeft   = "multiplayer.player.left"
)

const (
	StatusPending = iota
	StatusNormal
	StatusBanned
)
