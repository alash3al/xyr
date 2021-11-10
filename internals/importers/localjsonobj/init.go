package localjsonobj

import "github.com/alash3al/xyr/internals/kernel"

func init() {
	kernel.RegisterImporter("local+jsonobj", &Driver{})
}
