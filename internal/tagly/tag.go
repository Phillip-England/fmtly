package tagly

type Tag interface {
	TranspileToGo() (string, error)
	GetInfo() TagInfo
}
