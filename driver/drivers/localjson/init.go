package localjson

import "github.com/alash3al/xyr/driver"

func init() {
	driver.Register("local+json", &Driver{})
}
