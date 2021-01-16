package models

import "github.com/ypetya/fitfixer"

type Enhancer struct {
	// Find exercises to pass to the enhancer-function
	Filter IFilter
	// Enhancer function to call
	Function fitfixer.IEnhancer
	// Used by creating the new filename for enhanced data
	Name string
}
