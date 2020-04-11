package permission

type Permission struct {
	AllowUpload    bool
	AllowRename    bool
	AllowCreateDir bool
	AllowDelete    bool
	AllowMove      bool
	AllowDownload  bool
}

func (p Permission) HasOperation() bool {
	return p.AllowDownload || p.AllowDelete || p.AllowMove
}
