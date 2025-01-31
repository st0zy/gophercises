package adding

type Task struct {
	Name string
}

func NewTask(name string) Task {
	return Task{
		Name: name,
	}
}
