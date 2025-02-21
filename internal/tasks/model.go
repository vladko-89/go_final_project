package tasks

type Task struct {
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
	Title   string `json:"title"`
}

type TaskRow struct {
	Task
	ID string `json:"id"`
}

func NewTask(payload Task) *Task {
	return &Task{
		Date:    payload.Date,
		Comment: payload.Comment,
		Repeat:  payload.Repeat,
		Title:   payload.Title,
	}
}
