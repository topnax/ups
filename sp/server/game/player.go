package game

type Player struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	Ready        bool   `json:"ready"`
	Disconnected bool   `json:"disconnected"`
}
