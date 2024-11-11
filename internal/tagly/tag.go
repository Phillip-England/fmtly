package tagly

type Tag interface {
	AsStr() (string, error)
	GetInfo() TagInfo
	TranspileToGo() (string, error)
}
