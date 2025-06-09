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
