package db

type MockDb struct {
	titles map[int32]string
	about  map[int32]string
}

func NewMockDb() MockDb {
	return MockDb{
		titles: map[int32]string{
			1: "Some title nubmer 1",
			2: "This is the title nubmer 2",
			3: "Third title",
		},
		about: map[int32]string{
			1: "First about me section",
			2: "This is the second aboutme",
			3: "The third about me",
		},
	}
}

func (m MockDb) GetTitle(i int32) (string, error) {
	title, ok := m.titles[i]
	if !ok {
		return "title not found", nil
	}

	return title, nil
}

func (m MockDb) GetAboutMe(i int32) (string, error) {
	about, ok := m.about[i]
	if !ok {
		return "title not found", nil
	}

	return about, nil
}

func (m MockDb) Close() error {
	return nil
}
