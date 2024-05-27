package oraaq

type Message struct {
	ID      [16]byte
	Content string
}

func (m *Message) Raw() Message {
	return *m
}

func (m *Message) Text() string {
	return m.Content
}

func (m *Message) SetRaw(raw Message) {
	m.ID = raw.ID
	m.Content = raw.Content
}

func (m *Message) SetText(msg string) {
	m.Content = msg
}
