package handlers

import (
    "net/http"
    "strconv"
    "calc/mux-main"
    "calc/models"
    "calc/utils"
    "calc/storage"
)


type ExpressionHandler struct {
    Storage *storage.Storage
}

func (h *ExpressionHandler) GetExpressions(w http.ResponseWriter, r *http.Request) {
    expressions := h.Storage.GetExpressions()
    utils.RespondWithJSON(w, http.StatusOK, map[string][]models.Expression{"expression": expressions})
}

func (h *ExpressionHandler) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid expression ID")
        return
    }

    expression, err := h.Storage.GetExpressionByID(id)
    if err != nil {
        utils.RespondWithError(w, http.StatusNotFound, "Expression not found")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]models.Expression{"expression": expression})
}
