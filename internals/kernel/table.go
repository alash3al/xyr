package kernel

// Table represents a single internal storage table
type Table struct {
	Name    string   `hcl:"name,label"`
	DSN     string   `hcl:"dsn"`
	Loader  string   `hcl:"loader"`
	Columns []string `hcl:"columns"`

	ImporterInstance Importer
}
