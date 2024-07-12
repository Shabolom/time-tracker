package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"timeTracker/internal/model"
	"timeTracker/internal/utils"
	"timeTracker/pkg/logger"
)

// GetUsers godoc
// @Summary		Get all users
// @Tags			Users
// @Produce		json
// @Success		200	{object} []model.User
// @Param limit query string false "pagination limit"
// @Param page query string false "pagination page"
// @Param name query string false "filter name"
// @Param surname query string false "filter surname"
// @Param address query string false "filter address"
// @Param patronymic query string false "filter patronymic"
// @Router			/api/users [get]
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var filters model.UserFilter
	var pagination utils.Pagination

	filters.NameFilter = r.URL.Query().Get("name")
	filters.SurnameFilter = r.URL.Query().Get("surname")
	filters.AddressFilter = r.URL.Query().Get("address")
	filters.PatronymicFilter = r.URL.Query().Get("patronymic")

	if r.URL.Query().Get("limit") != "" {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		pagination.Limit = limit
	}

	if r.URL.Query().Get("page") != "" {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))

		if err != nil {
			h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		pagination.Page = page

	}

	pagination.Sort = r.URL.Query().Get("sort")

	users, err := h.Storage.GetUsers(ctx, filters, pagination)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	var userResponse []model.User
	for _, user := range users {
		userResponse = append(userResponse, model.User{
			Surname:    user.Surname,
			Name:       user.Name,
			Address:    user.Address,
			Patronymic: user.Patronymic,
			Base:       model.Base{ID: user.ID},
		})

	}

	err = h.Sender.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("Error when requesting /books")

		panic(err)
	}
}

// AddUser godoc
// @Summary		Add a specific user
// @Tags			Users
// @Produce		json
// @Accept			json
// @Param	user request	body		model.AddUserRequest	true	"user request"
// @Success		200	{object} model.User
// @Router			/api/users [post]
func (h *Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var passport model.AddUserRequest

	err := json.NewDecoder(r.Body).Decode(&passport)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(passport)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	match, err := regexp.MatchString("(\\d{2}\\s*\\d{2})\\s*(\\d{3,6})", passport.PassportNumber)

	if !match {
		h.Sender.JSON(w, http.StatusBadRequest, err)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", h.EnvBox.AuthService, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	pasportData := strings.Split(passport.PassportNumber, " ")

	passportSerie := pasportData[0]
	passportNumber := pasportData[1]

	q := req.URL.Query()
	q.Add("passportSerie", passportSerie)
	q.Add("passportNumber", passportNumber)
	req.URL.RawQuery = q.Encode()

	var resultAuth model.User
	var userEntity model.User
	userId := uuid.New()

	if h.EnvBox.Env == "DEV" {
		userEntity = model.User{
			Surname:    "Surname",
			Name:       "Name",
			Address:    "Address",
			Patronymic: "Patronymic",
			Base:       model.Base{ID: userId},
		}
	} else {
		resp, err := client.Do(req)

		if err != nil {
			h.Sender.JSON(w, http.StatusBadRequest, err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &resultAuth); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		userEntity = model.User{
			Surname:    resultAuth.Surname,
			Name:       resultAuth.Name,
			Address:    resultAuth.Address,
			Patronymic: resultAuth.Patronymic,
			Base:       model.Base{ID: userId},
		}
	}

	userResponse, err := h.Storage.AddUser(ctx, userEntity)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		panic(err)
	}
}

// UpdateUser godoc
// @Summary		Update a specific user
// @Tags			Users
// @Produce		json
// @Accept			json
// @Param	user request	body		model.UpdateUserRequest	true	"user request"
// @Success		200	{object} model.User
// @Router			/api/users [put]
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.UpdateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(user)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	userEntity := model.User{
		Surname:    user.Surname,
		Name:       user.Name,
		Address:    user.Address,
		Patronymic: user.Patronymic,
		Base:       model.Base{ID: user.ID},
	}

	userResponse, err := h.Storage.UpdateUser(ctx, userEntity)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		panic(err)
	}
}

// DeleteUser godoc
// @Summary		Delete a specific user
// @Tags			Users
// @Produce		json
// @Accept			json
// @Param	user request	body		model.DeleteUserRequest	true	"user request"
// @Success		200	{object} model.User
// @Router			/api/users [delete]
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.DeleteUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(user)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	ok, err := h.Storage.DeleteUser(ctx, user.ID)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, ok)
	if err != nil {
		panic(err)
	}
}
