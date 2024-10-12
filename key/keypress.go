package key

type KeyPress struct {
	Key      string `json:"key"`      // The key that was pressed (e.g., "a", "A", "Shift")
	ShiftKey bool   `json:"shiftKey"` // True if the Shift key was pressed
	CtrlKey  bool   `json:"ctrlKey"`  // True if the Ctrl key was pressed
	AltKey   bool   `json:"altKey"`   // True if the Alt key was pressed
	MetaKey  bool   `json:"metaKey"`  // True if the Meta (Command on Mac) key was pressed
}
