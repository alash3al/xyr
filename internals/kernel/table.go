package kernel

// Table represents a single internal storage table
type Table struct {
	Name         string   `hcl:"name,label"`
	ImporterName string   `hcl:"driver"`
	DSN          string   `hcl:"source"`
	Filter       string   `hcl:"filter"`
	Columns      []string `hcl:"columns"`

	ImporterInstance Importer
}
