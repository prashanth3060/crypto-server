package handler

import (
	"crypto-server/internal/httpserver"
	"crypto-server/internal/model"
	"crypto-server/internal/service"
	"crypto-server/internal/utils"
	"net/http"
)

type symbolRequest struct {
	BaseCurrency string `json:"base_currency"`
	PairCurrency string `json:"pair_currency"`
}

type SymbolHandler struct {
	symbol service.ISymbol
}

func NewSymbolHandler(symbol service.ISymbol) *SymbolHandler {
	return &SymbolHandler{
		symbol: symbol,
	}
}

func (s *SymbolHandler) GetRoutes() httpserver.Routes {
	return []httpserver.Route{
		{Verb: http.MethodGet, Path: "/symbol", Fn: s.GetAllSymbols},
		{Verb: http.MethodPost, Path: "/symbol", Fn: s.AddSymbol},
	}
}

func (s *SymbolHandler) AddSymbol(w http.ResponseWriter, r *http.Request) {

	symbolReq := &symbolRequest{}
	err := utils.DecodeRequest(r, symbolReq)
	if err != nil {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: err.Error(),
			Data:  nil,
		}, http.StatusBadRequest)
		return
	}

	if len(symbolReq.BaseCurrency) == 0 || len(symbolReq.PairCurrency) == 0 {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: "empty symbol/pair",
			Data:  nil,
		}, http.StatusBadRequest)
		return
	}

	err = s.symbol.AddSymbol(r.Context(), symbolReq.toSymbol())
	if err != nil {
		utils.EncodeResponseWithStatus(w, &Response{
			Error: err.Error(),
			Data:  nil,
		}, http.StatusInternalServerError)
		return
	}

	utils.EncodeResponseWithStatus(w, nil, http.StatusOK)

}

func (s *symbolRequest) toSymbol() *model.Symbol {
	return &model.Symbol{
		BaseCurrency: s.BaseCurrency,
		PairCurrency: s.PairCurrency,
		Symbol:       s.BaseCurrency + s.PairCurrency,
	}
}

func (s *SymbolHandler) GetAllSymbols(w http.ResponseWriter, r *http.Request) {
	symbols := s.symbol.GetAllSymbols(r.Context())

	utils.EncodeResponseWithStatus(w, &Response{
		Error: "",
		Data:  symbols,
	}, http.StatusOK)

}
