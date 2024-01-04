package shop_test

import (
	"fmt"
	"go-pbt/shop"
	"math/rand"
	"testing"
	"time"

	"pgregory.net/rapid"
)

type itemPrices struct {
	items         []string
	expectedPrice int
	prices        map[string]int
}

func itemPriceList(size int) *rapid.Generator[itemPrices] {
	return rapid.Custom(func(t *rapid.T) itemPrices {
		prices := priceList().Draw(t, "prices")
		items, expectedPrice := itemList(prices, size)
		return itemPrices{
			items:         items,
			expectedPrice: expectedPrice,
			prices:        prices,
		}
	})
}

func keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func random[K comparable, V any](m map[K]V, keys []K) (k K, v V) {
	if len(keys) == 1 {
		return keys[0], m[keys[0]]
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(keys) - 1)
	key := keys[i]
	return key, m[key]
}

func itemList(prices map[string]int, size int) (items []string, expectedPrice int) {
	items = make([]string, size)
	itemNames := keys(prices)
	for i := 0; i < size; i++ {
		item, price := random(prices, itemNames)
		items[i] = item
		expectedPrice += price
	}
	return items, expectedPrice
}

// key: itemName value: price
func priceList() *rapid.Generator[map[string]int] {
	return rapid.MapOfN(
		rapid.String().Filter(func(v string) bool { return v != "" }),
		rapid.Int(),
		1,
		100,
	)
}

func TestNoSpecial(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		size := rapid.IntRange(0, 30).Draw(t, "size")
		ip := itemPriceList(size).Draw(t, "ip")
		act, err := shop.Total(ip.items, ip.prices, []any{})
		if err != nil {
			t.Error(err)
		}

		if act != ip.expectedPrice {
			t.Errorf("expect total %d, but %d", ip.expectedPrice, act)
		}
	})
}

func TestPrices(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(priceList().Example(i))
	}
}
