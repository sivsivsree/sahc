package data

// HealthJobs is used to track the active running workers
type HealthJobs struct {
	// Running is the channel to stop the worker
	Running chan bool
}
