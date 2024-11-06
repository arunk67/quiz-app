package main

import "sync"

var (
	quizzes   = make(map[string]Quiz)
	results   = make(map[string]Result)
	quizMutex = &sync.Mutex{}
)
