package storage

type QueryParameter struct {
	Title           string
	Version         string
	MaintainerName  string
	MaintainerEmail string
	Company         string
	Website         string
	Source          string
	License         string
	Page			int
	PageSize		int
}

var ParameterNames = []string{"title", "version", "maintainerName", "maintainerEmail","company", "website", "source", "license"}
