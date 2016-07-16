package main

type WorkerGroup struct {
	num int
}

func NewWorkerGroup() *WorkerGroup {
	return &WorkerGroup{num: AppConfig.Workers}
}

func (self *WorkerGroup) Run() (tasks chan *VideoFile) {
	tasks = make(chan *VideoFile)
	for i := 0; i < self.num; i++ {
		logger.Debug("Spawned worker #", i+1)
		go Worker(tasks)
	}
	return tasks
}

func Worker(tasks chan *VideoFile) {
	for task := range tasks {
		task.RequestSubtitle()
	}
}
