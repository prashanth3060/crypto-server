package handler

import (
	"crypto-server/internal/httpserver"
	"crypto-server/internal/service"
	"crypto-server/internal/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type CurrencyHandler struct {
	currency service.ICurrency
}

func NewCurrencyHandler(event service.ICurrency) *CurrencyHandler {
	return &CurrencyHandler{
		currency: event,
	}
}

func (c *CurrencyHandler) GetRoutes() httpserver.Routes {
	return []httpserver.Route{
		{Verb: http.MethodGet, Path: "/currency/all", Fn: c.GetAllCurrency},
		{Verb: http.MethodGet, Path: "/currency/{symbol}", Fn: c.GetCurrency},
	}
}

func (c *CurrencyHandler) GetAllCurrency(w http.ResponseWriter, r *http.Request) {
	currencies, err := c.currency.GetAllCurrency(r.Context())
	if err != nil {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: err.Error(),
			Data:  nil,
		}, http.StatusInternalServerError)
		return
	}
	utils.EncodeResponseWithStatus(w, &Response{
		Error: "",
		Data:  currencies,
	}, http.StatusOK)
}

func (c *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := strings.ToUpper(vars["symbol"])
	currency, err := c.currency.GetCurrencyBySymbol(r.Context(), symbol)

	if err != nil {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: err.Error(),
			Data:  nil,
		}, http.StatusInternalServerError)
		return
	}

	if currency == nil {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: fmt.Sprintf("symbol %s not found", symbol),
			Data:  nil,
		}, http.StatusNotFound)
		return
	}
	utils.EncodeResponseWithStatus(w, &Response{
		Error: "",
		Data:  currency,
	}, http.StatusOK)
}
