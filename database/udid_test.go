package database

import (
	"testing"
)

func TestObjectID(t *testing.T) {
	pid = _genpid()

	// machineid 向左移位20
	machineid = _genmachineid() << 20
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
	println(ObjectID())
}
