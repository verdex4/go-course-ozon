package productservice

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	t.Run("successfully loads valid JSON", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		filePath := filepath.Join(dir, "products.json")

		content := `{
			"100": {"name": "Apple", "price": 10},
			"200": {"name": "Banana", "price": 20}
		}`
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		storage, err := NewStorage(filePath)
		require.NoError(t, err)
		require.NotNil(t, storage)

		prod, exists := storage.Get(100)
		require.True(t, exists)
		require.Equal(t, "Apple", prod.Name)
		require.Equal(t, uint32(10), prod.Price)

		prod, exists = storage.Get(200)
		require.True(t, exists)
		require.Equal(t, "Banana", prod.Name)
		require.Equal(t, uint32(20), prod.Price)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		filePath := filepath.Join(dir, "nonexistent.json")

		storage, err := NewStorage(filePath)
		require.Error(t, err)
		require.Nil(t, storage)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		filePath := filepath.Join(dir, "invalid.json")

		content := `{"100": {"name": "Apple", "price": 10` // Невалидный JSON
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		storage, err := NewStorage(filePath)
		require.Error(t, err)
		require.Nil(t, storage)
	})

	t.Run("returns error for invalid SKU format", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		filePath := filepath.Join(dir, "bad_sku.json")

		content := `{"not_a_number": {"name": "Apple", "price": 10}}`
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		storage, err := NewStorage(filePath)
		require.Error(t, err)
		require.Nil(t, storage)
	})
}

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns existing product", func(t *testing.T) {
		t.Parallel()

		storage := &Storage{
			products: map[int64]ProductModel{
				100: {Name: "Apple", Price: 10},
			},
		}

		prod, exists := storage.Get(100)
		require.True(t, exists)
		require.Equal(t, "Apple", prod.Name)
		require.Equal(t, uint32(10), prod.Price)
	})

	t.Run("returns false for non-existent product", func(t *testing.T) {
		t.Parallel()

		storage := &Storage{
			products: map[int64]ProductModel{},
		}

		prod, exists := storage.Get(999)
		require.False(t, exists)
		require.Equal(t, ProductModel{}, prod)
	})
}
