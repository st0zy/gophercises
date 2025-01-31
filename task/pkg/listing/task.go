package listing

type Task struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
