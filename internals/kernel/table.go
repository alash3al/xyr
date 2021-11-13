package kernel

// Table represents a single internal storage table
type Table struct {
	Name         string   `hcl:"name,label"`
	ImporterName string   `hcl:"driver"`
	DSN          string   `hcl:"source"`
	Loader       string   `hcl:"loader"`
	Columns      []string `hcl:"columns"`

	ImporterInstance Importer
}
