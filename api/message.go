package api

type Message[R any] interface {
	Raw() R
	Text() string
	SetRaw(R)
	SetText(string)
}
