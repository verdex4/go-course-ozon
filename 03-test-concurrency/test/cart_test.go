//go:build api

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

func TestCartService_Integration(t *testing.T) {
	userID := time.Now().UnixNano() % 1000000

	t.Run("GetCart returns empty cart for new user", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s/user/%d/cart", cartServiceAddr, userID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var cart CartResponse
		err = json.NewDecoder(resp.Body).Decode(&cart)
		require.NoError(t, err)
		require.Empty(t, cart.Items)
		require.Equal(t, uint32(0), cart.TotalPrice)
	})

	t.Run("AddProduct and GetCart happy path", func(t *testing.T) {
		t.Parallel()

		// SKU 1076963 есть в data.json
		skuID := int64(1076963)
		count := uint32(3)

		url := fmt.Sprintf("%s/user/%d/cart/%d", cartServiceAddr, userID, skuID)
		body := bytes.NewBufferString(fmt.Sprintf(`{"count": %d}`, count))
		req, err := http.NewRequest(http.MethodPost, url, body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		getURL := fmt.Sprintf("%s/user/%d/cart", cartServiceAddr, userID)
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		var cart CartResponse
		err = json.NewDecoder(getResp.Body).Decode(&cart)
		require.NoError(t, err)
		require.Len(t, cart.Items, 1)
		require.Equal(t, skuID, cart.Items[0].Sku)
		require.Equal(t, count, cart.Items[0].Count)
		require.Greater(t, cart.Items[0].Price, uint32(0))
		require.Greater(t, cart.TotalPrice, uint32(0))
	})

	t.Run("AddProduct returns 412 for non-existent SKU", func(t *testing.T) {
		t.Parallel()

		nonExistentSKU := int64(999999999)
		url := fmt.Sprintf("%s/user/%d/cart/%d", cartServiceAddr, userID, nonExistentSKU)
		body := bytes.NewBufferString(`{"count": 1}`)

		resp, err := http.Post(url, "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusPreconditionFailed, resp.StatusCode)
	})

	t.Run("GetCart returns 400 for invalid user_id", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s/user/abc/cart", cartServiceAddr)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("AddProduct returns 400 for invalid request body", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s/user/%d/cart/1", cartServiceAddr, userID)
		body := bytes.NewBufferString(`invalid json`)

		resp, err := http.Post(url, "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("AddProduct returns 400 for count <= 0", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s/user/%d/cart/1", cartServiceAddr, userID)
		body := bytes.NewBufferString(`{"count": 0}`)

		resp, err := http.Post(url, "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
