package cart

import (
	"fmt"

	"github.com/olliekm/gorestapi/types"
)

func getCartIDs(cartItems []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(cartItems))
	for i, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product ID %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOder(products []types.Product, cartItems []types.CartItem, userID int) (int, float64, error) {
	// create a map of products for easier access
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	// check if all products are available
	if err := checkProductAvailability(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	// reduce the quantity of products in the store
	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProductStock(product)
	}

	// create order record
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // could fetch address from a user addresses table
	})
	if err != nil {
		return 0, 0, err
	}

	// create order the items records
	for _, item := range cartItems {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkProductAvailability(cartItems []types.CartItem, productMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, exists := productMap[item.ProductID]
		if !exists {
			return fmt.Errorf("product ID %d does not exist", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("not enough stock for product ID %d", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, productMap map[int]types.Product) float64 {
	total := 0.0
	for _, item := range cartItems {
		product := productMap[item.ProductID]
		total += float64(item.Quantity) * product.Price
	}
	return total
}
