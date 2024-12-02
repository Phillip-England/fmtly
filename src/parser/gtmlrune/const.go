package gtmlrune

const (
	KeyRuneProp = "$prop"
	KeyRuneSlot = "$slot"
	KeyRuneVal  = "$val"
	KeyRunePipe = "$pipe"
	KeyRuneMd   = "$md"
)

const (
	KeyLocationAttribute = "KEYLOCATIONATTRIBUTE"
	KeyLocationElsewhere = "KEYLOCATIONELSEWHERE"
)

func GetRuneNames() []string {
	return []string{KeyRuneProp, KeyRuneSlot, KeyRuneVal, KeyRunePipe, KeyRuneMd}
}
