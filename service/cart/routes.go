package cart

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olliekm/gorestapi/types"
	"github.com/olliekm/gorestapi/utils"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handleCheckout).Methods("POST")

}

func (h *Handler) handleCheckout(w http.ResponseWriter, req *http.Request) {
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(req, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get products
	productIDs, err := getCartIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
}
