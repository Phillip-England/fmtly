package tag

type Tag interface {
	TranspileToGo() (string, error)
	GetInfo() TagInfo
}
