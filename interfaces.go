package CarCatalog

type Repo interface {
	GetAll(queryParamsMap map[string]string, offset int) ([]Car, error)
	DeleteByID(regNum string) error
	ChangeByID(car Car) error
	AddNew(cars []Car) error
}
