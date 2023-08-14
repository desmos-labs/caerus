package files

type Database interface {
	SaveMediaHash(imageUrl string, hash string) error
}
