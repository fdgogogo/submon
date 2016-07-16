package main

type TaskQueue struct {
	num int
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{num: AppConfig.Workers}
}

func (self *TaskQueue) Run() (tasks chan *VideoFile) {
	tasks = make(chan *VideoFile)
	for i := 0; i < self.num; i++ {
		logger.Debugf("Spawned worker #%d", i+1)
		go Worker(tasks)
	}
	return tasks
}

func Worker(tasks chan *VideoFile) {
	for task := range tasks {
		task.RequestSubtitle()
	}
}
