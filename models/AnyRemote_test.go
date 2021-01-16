package models

import "testing"

var (
	r   = AnyRemote{Name: "garmin", Target: GarminConnect}
	r1  = AnyRemote{Name: "garmin", Target: GarminConnect}
	r2  = AnyRemote{Name: "garmin", Target: Directory}
	r3  = AnyRemote{Name: "different", Target: GarminConnect}
	rl1 = AnyRemote{Target: LocalDB}
	rl2 = AnyRemote{Name: "different", Target: LocalDB}
)

func TestAREqualsWhenNameAndTargetEqualsOrItIsALocalDB(t *testing.T) {
	if !r.equals(r1) {
		t.Error("remotes should be equal by name and target being the same")
	}
	if r.equals(r2) || r.equals(r3) {
		t.Error("remotes should be equal by name and target being the same")
	}
	if !rl1.equals(rl2) {
		t.Error("only one localdb is posible!")
	}
}
