package jobs

type WorkerPool interface {
	StopAndWait()
	Stopped() bool

	Submit(func())
}
