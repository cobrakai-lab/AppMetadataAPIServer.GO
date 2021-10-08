package Model

type AppMetadata struct {
	Title       string       `json:"title"`
	Version     string       `json:"version"`
	Maintainers []Maintainer `json:"maintainers"`
	Company     string       `json:"company"`
	Website     string       `json:"website"`
	Source      string       `json:"source"`
	License     string       `json:"license"`
	Description string       `json:"description"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
