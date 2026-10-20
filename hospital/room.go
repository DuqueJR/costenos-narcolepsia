package hospital

// Room is a hospital room with a fixed number of beds. Its fields are
// unexported, so patients enter and leave only through its methods.
type Room struct {
	number    int
	capacity  int
	state     RoomState
	occupants []*Patient
}
