package jsondir

import "github.com/alash3al/xyr/internals/kernel"

func init() {
	kernel.RegisterImporter("jsondir", &Driver{})
}
