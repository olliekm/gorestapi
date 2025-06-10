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

func (h *Handler) createOder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, p := range ps {
		productMap[p.ID] = p
	}

	// Check that all products are available
	if err := checkProductAvailability(items, productMap); err != nil {
		return 0, 0, nil
	}

	// Calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// TODO: Atomize
	// Reduce stock for each product
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProductStock(product)
	}

	// Create order
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Default Address", // This should be replaced with actual address handling
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
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
