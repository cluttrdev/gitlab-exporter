package types

type JobArtifact struct {
	Job JobReference

	FileType     string
	Name         string
	DownloadPath string
}
