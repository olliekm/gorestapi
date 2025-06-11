package cart

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olliekm/gorestapi/service/auth"
	"github.com/olliekm/gorestapi/types"
	"github.com/olliekm/gorestapi/utils"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods("POST")

}

func (h *Handler) handleCheckout(w http.ResponseWriter, req *http.Request) {
	userID := auth.GetUserIDFromContext(req.Context())
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
		fmt.Println("1")
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		fmt.Println("2")

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOder(ps, cart.Items, userID)
	if err != nil {
		fmt.Println("3")

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})

}
