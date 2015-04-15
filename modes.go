package main

type Mode interface {
	Major() bool
	Name() string
}

type Modes []Mode
