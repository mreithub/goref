package goref

// singleton GoRef instance
var instance = NewGoRef()

// Clone -- Returns a Snapshot of the GoRef  (synchronously)
func Clone() *Snapshot {
	return instance.Clone()
}

// Ref -- References an instance of 'key' (in singleton mode)
func Ref(key string) Instance {
	return instance.Ref(key)
}

// Reset -- replaces the GoRef singleton instance
// Note: This function is NOT synchronized and you might end up losing some data.
//   But since 'losing data' is the idea of this function, there's really no
//   downside (at least none I can think of).
//   If that's a problem for you, create your own GoRef instance and ignore the
//   singleton.
func Reset() {
	instance = NewGoRef()
}

// Snapshots -- Linked list of the singleton instance's Snapshots (in reverse order)
func Snapshots() *Snapshot {
	return instance.Snapshots()
}

// TakeSnapshot -- singleton version of NewGoRef().TakeSnapshot()
func TakeSnapshot(name string) *Snapshot {
	return instance.TakeSnapshot(name)
}
