package messages

// Server <-> Client (Sync)
//procm:use=derive_binary
type GetChestName struct {
	ChestID int16
	ChestX  int16
	ChestY  int16
	Name    string
}
