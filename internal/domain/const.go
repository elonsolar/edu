package domain

type SubjectCategary int

const (
	PAINTING SubjectCategary = iota
	HANDWRITING
)

func (s SubjectCategary) String() string {
	return []string{"绘画", "书法"}[s]
}

func (w SubjectCategary) Original() int {
	return int(w)
}
