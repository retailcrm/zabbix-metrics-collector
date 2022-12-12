package metrics

func StoppableFinalizer(s Stoppable) {
	_ = s.Stop()
}
