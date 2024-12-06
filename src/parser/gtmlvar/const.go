package gtmlvar

const (
	KeyVarGoFor         = "VARGOFOR"
	KeyVarGoIf          = "VARGOIF"
	KeyVarGoElse        = "VARGOELSE"
	KeyVarGoPlaceholder = "VARGOPLACEHOLDER"
	KeyVarGoSlot        = "VARGOSLOT"
	KeyVarGoMd          = "VARGOMD"
)

func GetFullVarList() []string {
	return []string{KeyVarGoFor, KeyVarGoIf, KeyVarGoElse, KeyVarGoPlaceholder, KeyVarGoSlot, KeyVarGoMd}
}
