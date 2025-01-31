package adding

type Task struct {
	Id uint64
}

func ToTask(id uint64) Task {
	return Task{
		Id: id,
	}
}
