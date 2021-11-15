package s3jsondir

import "github.com/alash3al/xyr/internals/kernel"

func init() {
	kernel.RegisterImporter("s3jsondir", &Driver{})
}
