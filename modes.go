package main

import (
	"github.com/kiasaki/ry/frontends"
)

type Mode interface {
	Major() bool
	Name() string
	HandleKey(*Editor, frontends.Event)
}

type Modes []Mode
