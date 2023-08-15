package files

type Database interface {
	SaveMediaHash(fileName string, hash string) error
}
