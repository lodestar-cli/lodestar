package remote

type LodestarFile interface {
	Print()
	Output() error
	GetByteContent() []byte
	GetStringContent() string
}