package db

import (
	"testing"

	"github.com/Logmemo/price-watch/parse"
)

func TestParseProduct(t *testing.T) {

	db := Connect()
	_ = db

}

func TestAddProduct(t *testing.T) {

	product := parse.ParseProduct("https://www.ozon.ru/product/1880425410")
	AddProduct(product)

}

func TestAddProduct2(t *testing.T) {

	product := parse.ParseProduct("https://www.ozon.ru/product/4755490404")
	AddProduct(product)

}
