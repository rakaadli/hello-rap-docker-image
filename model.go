package main

import "time"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Birth     time.Time
}
