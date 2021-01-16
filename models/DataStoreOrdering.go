package models

// Len is the number of elements in the collection.
func (ds DataStore) Len() int {
	return len(ds.Exercises)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (ds DataStore) Less(i, j int) bool {
	return ds.Exercises[i].Meta.Created < ds.Exercises[j].Meta.Created
}

// Swap swaps the elements with indexes i and j.
func (ds DataStore) Swap(i, j int) {
	ds.Exercises[i], ds.Exercises[j] = ds.Exercises[j], ds.Exercises[i]
}
