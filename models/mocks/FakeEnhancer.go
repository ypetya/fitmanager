package mocks

type FakeEnhancer struct{}

func (FakeEnhancer) Enhance(targetFile string, toEnhance string, with string) {
	LogCall("Enhance:" + targetFile + "," + toEnhance + "," + with)
}
