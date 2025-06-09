package product

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olliekm/gorestapi/types"
	"github.com/olliekm/gorestapi/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/create-product", h.handleCreateProduct).Methods("POST")

}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// get products from store
	products, err := h.store.GetProducts() // Assuming 0 means no specific product ID filter
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// write products as JSON response
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// parse JSON payload
	var product types.ProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate product
	if err := utils.Validate.Struct(product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// create product in store
	if err := h.store.CreateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}
