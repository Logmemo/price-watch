package parse

import (
	"log"
	"strconv"
	"strings"

	"github.com/mxschmitt/playwright-go"
)

type Product struct {
	Title         string
	BasicPrice    int
	DiscountPrice int
	Image         string
	Url           string
}

func ParseProduct(url string) Product {

	var product Product
	product.Url = url

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Не удалось запустить Playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--disable-features=IsolateOrigins,site-per-process",
			"--no-sandbox",
			"--disable-web-security",
		},
	})
	if err != nil {
		log.Fatalf("Не удалось запустить Chromium: %v", err)
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent:         playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.0"),
		Locale:            playwright.String("ru-RU"),
		TimezoneId:        playwright.String("Europe/Moscow"),
		Viewport:          &playwright.Size{Width: 1920, Height: 1080},
		Screen:            &playwright.Size{Width: 1920, Height: 1080},
		DeviceScaleFactor: playwright.Float(1),
		IsMobile:          playwright.Bool(false),
		HasTouch:          playwright.Bool(false),
		ColorScheme:       playwright.ColorSchemeLight,
		AcceptDownloads:   playwright.Bool(false),
	})

	if err != nil {
		log.Fatalf("Не удалось создать контекст: %v", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Не удалось создать страницу: %v", err)
	}

	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("Ошибка перехода: %v", err)
	}

	product.Title, err = page.Locator("head title").TextContent()
	if err != nil {
		log.Fatalf("Не удалось найти h1: %v", err)
	}

	discountPrice, err := page.Locator("div[data-widget='webPrice'] span.tsHeadline600Large").TextContent()
	if err != nil {
		log.Fatalf("Не удалось найти цену со скидкой: %v", err)
	}
	product.DiscountPrice = priceCleaner(discountPrice)

	basicPrice, err := page.Locator("div[data-widget='webPrice'] span.tsHeadline500Medium").TextContent()
	if err != nil {
		log.Fatalf("Не удалось найти обычную цену: %v", err)
	}
	product.BasicPrice = priceCleaner(basicPrice)

	product.Image, err = page.Locator("head link[rel='preload'][as='image']").GetAttribute("href")
	if err != nil {
		log.Fatalf("Не удалось найти изображение товара: %v", err)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("Не удалось закрыть Chromium: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("Не удалось остановить Playwright: %v", err)
	}

	return product
}

func priceCleaner(price string) int {
	var cleanPrice int
	var err error
	var currency rune

	for _, r := range price {
		currency = r
	}

	switch currency {
	case '₽':
		replacer := strings.NewReplacer(
			"\u2009", "",
			"₽", "",
		)
		clean := replacer.Replace(price)
		cleanPrice, err = strconv.Atoi(clean)
		if err != nil {
			log.Fatalf("ошибка при преобразовании типа цены(рубль): %v", err)
		}
		return cleanPrice
	case '$':
		replacer := strings.NewReplacer(
			"\u2009", "",
			"$", "",
		)
		clean := replacer.Replace(price)
		cleanPrice, err = strconv.Atoi(clean)
		if err != nil {
			log.Fatalf("ошибка при преобразовании типа цены(доллар): %v", err)
		}
		return cleanPrice
	case '₸':
		replacer := strings.NewReplacer(
			"\u2009", "",
			"₸", "",
		)
		clean := replacer.Replace(price)
		cleanPrice, err = strconv.Atoi(clean)
		if err != nil {
			log.Fatalf("ошибка при преобразовании типа цены(тэнге): %v", err)
		}
		return cleanPrice
	default:
		var sb strings.Builder
		for _, r := range price {
			if r >= '0' && r <= '9' {
				sb.WriteRune(r)
			}
		}
		if n, err := strconv.Atoi(sb.String()); err == nil {
			return n
		}
		return 0
	}
}
