package models

type File struct {
	Path     string `json:"path,omitempty"`
	LastSeen int64  `json:"lastSeen,omitempty"`
}

func (f *File) HasName() bool {
	return len(f.Path) > 0
}

func (f *File) AsDir() string {
	ix := len(f.Path) - 1
	if ix >= 0 && f.Path[ix] != '/' {
		return f.Path + "/"
	}
	return f.Path
}
