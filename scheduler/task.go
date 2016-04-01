package scheduler

type TaskDetail struct {
	script string
	fun    func(name string)
}

func NewTask(script string, fun func(name string)) *TaskDetail {
	return &TaskDetail{script: script, fun: fun}
}

func (j *TaskDetail) Run() {
	j.fun(j.script)
}
