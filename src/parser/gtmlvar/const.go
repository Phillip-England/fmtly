package gtmlvar

const (
	KeyVarGoFor         = "VARGOFOR"
	KeyVarGoIf          = "VARGOIF"
	KeyVarGoElse        = "VARGOELSE"
	KeyVarGoPlaceholder = "VARGOPLACEHOLDER"
	KeyVarGoSlot        = "VARGOSLOT"
)

func GetFullVarList() []string {
	return []string{KeyVarGoFor, KeyVarGoIf, KeyVarGoElse, KeyVarGoPlaceholder, KeyVarGoSlot}
}
