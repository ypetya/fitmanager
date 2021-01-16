package mocks

type FakeConnector struct {
	FakeEntities []string
}

// TODO implement fake methods for import/ update:

func (f FakeConnector) FetchDiff(knownExternalIds []string) []string {
	LogCall("FetchDiff")
	var ret []string
outer:
	for i, v := range f.FakeEntities {
		for _, k := range knownExternalIds {
			if k == v {
				continue outer
			}
		}
		ret = append(ret, f.FakeEntities[i])
	}
	return ret
}

func (f FakeConnector) SetTarget(dirName string) {
	LogCall("SetTarget:" + dirName)
}

func (f FakeConnector) SetSource(dirName string) {
	LogCall("SetSource:" + dirName)
}

func CreateFakeConnector(entities []string) FakeConnector {
	FakeCalls = []string{}
	return FakeConnector{FakeEntities: entities}
}

func (f FakeConnector) Import(targetDir string, externalId string) (string, error) {
	LogCall("Import:" + externalId)
	return externalId, nil
}

func (f FakeConnector) Export(source string, id string) (string, error) {
	LogCall("Export:" + source + id)
	return id, nil
}

func (f FakeConnector) Update(localDir string, localId string, remoteId string) error {
	LogCall("Update:" + localDir + localId + remoteId)
	return nil
}
