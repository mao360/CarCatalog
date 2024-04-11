package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mao360/CarCatalog"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var PagSize = 5

type Handler struct {
	repo   CarCatalog.Repo
	logger *logrus.Logger
}

func NewHandler(repo CarCatalog.Repo, log *logrus.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: log,
	}
}

// @Summary Get All cars with filtration
// @Description Get All cars with filtration
// @ID get-all
// @Accept  json
// @Produce  json
// @Success 200 {object} []CarCatalog.Car
// @Failure 500 {object} error
// @Router /cars [get]
func (h *Handler) GetAll(c echo.Context) error {
	h.logger.Infof("handlers.go, GetAll started")

	queryParamsMap := make(map[string]string)
	queryParamsMap["mark"] = c.QueryParam("mark")
	queryParamsMap["model"] = c.QueryParam("model")
	queryParamsMap["year"] = c.QueryParam("year")
	queryParamsMap["reg_num"] = c.QueryParam("reg_num")
	queryParamsMap["name"] = c.QueryParam("name")
	queryParamsMap["surname"] = c.QueryParam("surname")
	queryParamsMap["patronymic"] = c.QueryParam("patronymic")

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.logger.Errorf("handlers.go, GetAll, strconv.Atoi: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, GetAll, strconv.Atoi: %w", err)
	}
	offset := (page - 1) * PagSize
	result, err := h.repo.GetAll(queryParamsMap, offset)
	if err != nil {
		h.logger.Errorf("handlers.go, GetAll, repo.GetAll: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, repo.GetAll, %w", err)
	}

	h.logger.Infof("handlers.go, GetAll successfully finished")
	return c.JSON(http.StatusOK, &result)
}

// @Summary Delete car from list by regNum
// @Description Delete car from list by regNum
// @ID delete-by-id
// @Success 200 no content
// @Failure 500 {object} error
// @Router /cars [delete]
func (h *Handler) DeleteByID(c echo.Context) error {

	h.logger.Infof("handlers.go, DeleteByID started")

	regNum := c.QueryParam("reg_num")

	err := h.repo.DeleteByID(regNum)
	if err != nil {
		h.logger.Errorf("handlers.go, DeleteByID, repo.DeleteByID: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, DeleteByID, repo.DeleteByID: %w", err)
	}

	c.NoContent(http.StatusOK)
	h.logger.Infof("handlers.go, DeleteByID successfully finished")
	return nil
}

// @Summary Change record by ID
// @Description Change record by ID
// @ID change-by-id
// @Success 200 no content
// @Failure 500 {object} error
// @Router /cars [put]
func (h *Handler) ChangeByID(c echo.Context) error {
	h.logger.Infof("handlers.go, ChangeByID started")

	regNum := c.QueryParam("reg_num")
	mark := c.QueryParam("mark")
	model := c.QueryParam("model")
	strYear := c.QueryParam("year")

	var year int
	var err error
	if strYear != "" {
		year, err = strconv.Atoi(strYear)
		if err != nil {
			h.logger.Errorf("handlers.go, ChangeByID, strconv.Atoi: %s", err.Error())
			c.NoContent(http.StatusInternalServerError)
			return fmt.Errorf("handlers.go, ChangeByID, strconv.Atoi: %w", err)
		}
	}

	car := CarCatalog.Car{
		RegNum: regNum,
		Mark:   mark,
		Model:  model,
		Year:   year,
	}

	err = h.repo.ChangeByID(car)
	if err != nil {
		h.logger.Errorf("handlers.go, ChangeByID, repo.ChangeByID: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, ChangeByID, repo.ChangeByID: %w", err)
	}
	c.NoContent(http.StatusOK)
	h.logger.Infof("handlers.go, ChangeByID successfully finished")
	return nil
}

// @Summary Create new car record
// @Description Create new car record
// @ID add-new
// @Success 200 no content
// @Failure 500 {object} error
// @Router /cars [post]
func (h *Handler) AddNew(c echo.Context) error {
	h.logger.Infof("handlers.go, AddNew started")

	data, err := ioutil.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		h.logger.Errorf("handlers.go, AddNew, ioutil.ReadAll: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, AddNew, ioutil.ReadAll: %w", err)
	}
	var RegNums struct {
		RegNumsIn []string `json:"regNums"`
	}
	var Temp struct {
		RegNumsIn json.RawMessage `json:"regNums"`
	}
	err = json.Unmarshal(data, &Temp)
	if err != nil {
		h.logger.Errorf("handlers.go, AddNew, json.Unmarshal: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, AddNew, json.Unmarshal: %w", err)
	}

	err = json.Unmarshal(Temp.RegNumsIn, &RegNums.RegNumsIn)
	if err != nil {
		h.logger.Errorf("handlers.go, AddNew, json.Unmarshal: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, AddNew, json.Unmarshal: %w", err)
	}

	client := &http.Client{}
	cars := make([]CarCatalog.Car, 0)
	car := new(CarCatalog.Car)
	for _, val := range RegNums.RegNumsIn {
		url := fmt.Sprintf("%s/info?name=%s", os.Getenv("EXTERNAL_SERVICE_DOMAIN"), val)
		resp, err := client.Get(url)
		defer resp.Body.Close()
		if err != nil {
			h.logger.Errorf("handlers.go, AddNew, clientGet: %s", err.Error())
			c.NoContent(http.StatusInternalServerError)
			return fmt.Errorf("handlers.go, AddNew, clientGet: %w", err)
		}
		externalServBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("handlers.go, AddNew, ioutil.ReadAll: %s", err.Error())
			c.NoContent(http.StatusInternalServerError)
			return fmt.Errorf("handlers.go, AddNew, ioutil.ReadAll: %w", err)
		}

		err = json.Unmarshal(externalServBody, car)
		if err != nil {
			h.logger.Errorf("handlers.go, AddNew, json.Unmarshal: %s", err.Error())
			c.NoContent(http.StatusInternalServerError)
			return fmt.Errorf("handlers.go, AddNew, json.Unmarshal: %w", err)
		}

		cars = append(cars, *car)
	}

	err = h.repo.AddNew(cars)
	if err != nil {
		h.logger.Errorf("handlers.go, AddNew, repo.AddNew: %s", err.Error())
		c.NoContent(http.StatusInternalServerError)
		return fmt.Errorf("handlers.go, AddNew, repo.AddNew: %w", err)
	}

	c.NoContent(http.StatusCreated)
	h.logger.Infof("handlers.go: AddNew successfully finished")
	return nil
}
