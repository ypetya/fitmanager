package mocks

import "github.com/ypetya/fitmanager/models"

type FakeFilter struct {
	Calls  []string
	Params [][]string
}

func (ff *FakeFilter) Remote(remote string, condition []string) {
	ff.Calls = append(ff.Calls, "Remote")
	condi := []string{}
	condi = append(condi, remote)
	condi = append(condi, condition...)
	ff.Params = append(ff.Params, condi)
}

func (ff *FakeFilter) Filter(e []models.Exercise) []models.Exercise {
	ff.Calls = append(ff.Calls, "Filter")
	return []models.Exercise{}
}
