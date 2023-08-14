package files

import (
	"github.com/desmos-labs/caerus/routes/base"
)

type Database interface {
	base.Database
	SaveMediaHash(imageUrl string, hash string) error
}
