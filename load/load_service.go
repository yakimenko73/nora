package load

type LoadService interface {
	Start()
	AddJob(jobFunc func() error) error
	SetLoadTime(loadTime float32)
}
