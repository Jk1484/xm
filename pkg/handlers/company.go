package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"xm/pkg/repositories/company"
	"xm/pkg/services/utils"
)

func (h *handlers) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	var c company.Company
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	err = validate(c)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	err = h.companyService.Create(c)
	if err != nil {
		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		h.logger.Logger().Error(err)
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), "created")
}

func (h *handlers) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	idS := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "bad id")
		return
	}

	c, err := h.companyService.GetByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			apiResp.Set(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
			return
		}

		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		h.logger.Logger().Error(err)
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), c)
}

func (h *handlers) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	var f company.Filters
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	c, err := h.companyService.GetAll(f)
	if err != nil {
		if err == utils.ErrNotFound {
			apiResp.Set(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
			return
		}

		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		h.logger.Logger().Error(err)
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), c)
}

func (h *handlers) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	var c company.Company
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	err = h.companyService.Update(c)
	if err != nil {
		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		h.logger.Logger().Error(err)
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), "updated")
}

func (h *handlers) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	var apiResp ApiResp
	defer apiResp.Respond(w)

	idS := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		apiResp.Set(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "bad id")
		return
	}

	err = h.companyService.DeleteByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			apiResp.Set(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
			return
		}

		apiResp.Set(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		h.logger.Logger().Error(err)
		return
	}

	apiResp.Set(http.StatusOK, http.StatusText(http.StatusOK), "deleted")
}

func validate(c company.Company) error {
	if c.Name == "" {
		return errors.New("no name provided")
	}

	if c.Code == "" {
		return errors.New("no code provided")
	}

	if c.Country == "" {
		return errors.New("no country provided")
	}

	if c.Website == "" {
		return errors.New("no website provided")
	}

	if c.Phone == "" {
		return errors.New("no phone provided")
	}

	return nil
}
