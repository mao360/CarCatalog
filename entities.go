package CarCatalog

type Car struct {
	CarID  string `json:"-"`
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  Human  `json:"owner"`
}

type Human struct {
	HumanID    string `json:"-"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}
