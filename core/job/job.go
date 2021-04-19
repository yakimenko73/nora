package job

type Job interface {
	Complete() error
	Cancel()
}
