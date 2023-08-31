package applications

type Database interface {
	IsUserAdminOfApp(desmosAddress string, appID string) (bool, error)
	DeleteApp(appID string) error
}
