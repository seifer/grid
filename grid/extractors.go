package grid

import (
	"encoding/json"
	"strings"

	html "github.com/antchfx/xquery/html"
	xhtml "golang.org/x/net/html"
)

func extractTitle(node *xhtml.Node) string {
	nodeTitle := html.FindOne(node, "//span[@id='productTitle']")
	if nodeTitle != nil && nodeTitle.FirstChild != nil {
		return strings.TrimSpace(nodeTitle.FirstChild.Data)
	}

	return ""
}

func extractImage(node *xhtml.Node) string {
	nodeImage := html.FindOne(node, "//img[@id='landingImage' or @id='imgBlkFront']")

	if nodeImage != nil {
		for i := range nodeImage.Attr {
			if nodeImage.Attr[i].Key == "data-a-dynamic-image" {
				attrs := map[string]interface{}{}

				if err := json.Unmarshal([]byte(nodeImage.Attr[i].Val), &attrs); err != nil {
					return ""
				}

				for key := range attrs {
					return strings.TrimSpace(key)
				}
			}
		}
	}

	return ""
}

func extractPrice(node *xhtml.Node) string {
	nodePrice := html.FindOne(node, "//span[@id='priceblock_ourprice']")
	if nodePrice != nil && nodePrice.FirstChild != nil {
		return strings.TrimSpace(nodePrice.FirstChild.Data)
	}

	nodePrice = html.FindOne(node, "//span[contains(@class,'offer-price')]")
	if nodePrice != nil && nodePrice.FirstChild != nil {
		return strings.TrimSpace(nodePrice.FirstChild.Data)
	}

	return ""
}

func extractIsInStock(node *xhtml.Node) bool {
	return html.FindOne(node, "//div[@id='availability']") != nil
}
