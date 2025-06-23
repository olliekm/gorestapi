package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/olliekm/gorestapi/types"
)

type Store struct {
	db          *sql.DB
	redisClient *redis.Client
	cacheTTL    time.Duration
}

func NewStore(db *sql.DB, redisClient *redis.Client) *Store {
	return &Store{
		db:          db,
		redisClient: redisClient,
		cacheTTL:    5 * time.Minute, // Set cache TTL to 5 minutes
	}
}

func (s *Store) GetProducts(ctx context.Context) ([]*types.Product, error) {
	// Chaching
	const cacheKey = "products:all"

	if data, err := s.redisClient.Get(ctx, cacheKey).Bytes(); err == nil {
		// If data is found in Redis cache, unmarshal it into products
		var products []*types.Product
		if err := json.Unmarshal(data, &products); err != nil {
			return products, nil
		}
	}
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}
	// 3) Marshal & set cache (best effort)
	if data, err := json.Marshal(products); err == nil {
		s.redisClient.Set(ctx, cacheKey, data, s.cacheTTL)
	}

	return products, nil
}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) CreateProduct(product types.ProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.ImageURL, product.Price, product.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert productIDs to []interface{}
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProductStock(product types.Product) error {
	fmt.Print("Updating product stock: ", product.ID, "\n")
	_, err := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Price, product.ImageURL, product.Description, product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil
}
