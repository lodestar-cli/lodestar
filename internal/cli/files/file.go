package files

type LodestarFile interface {
	Print()
	Output() error
	GetByteContent() []byte
	GetStringContent() string
}