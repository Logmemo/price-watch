package parse

import "testing"

func TestParseProduct(t *testing.T) {

	url := "https://www.ozon.ru/product/1880425410"
	expected := Product{"Go: идиомы и паттерны проектирования, 2-е издание. Язык программирования Golang | Боднер Джон купить на OZON по низкой цене (1880425410)", 2221, 1999, "https://ir.ozone.ru/s3/multimedia-1-o/wc1000/9122389692.jpg", "https://www.ozon.ru/product/1880425410"}

	result := ParseProduct(url)

	if expected != result {
		t.Errorf("\nТест провален, получено: %#v, ожидалось: %#v.", result, expected)
	}
}

func TestPriceCleaner(t *testing.T) {

	price := "1 999 ₽"
	expected := 1999

	result := priceCleaner(price)

	if expected != result {
		t.Errorf("\nТест провален, получено: %#v, ожидалось: %#v.", result, expected)
	}
}
