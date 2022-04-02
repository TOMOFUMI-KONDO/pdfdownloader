package pdfdownloader

import (
	"golang.org/x/net/html"
)

func sessionPageLinks(doc *html.Node) []string {
	htmlElm := nthChild(doc, 2)
	bodyElm := findByData(htmlElm, "body")
	pageBodyElm := findByAttr(bodyElm, "id", "pagebody")
	contentsElm := findByAttr(pageBodyElm, "id", "contents")
	menuElm := findByAttr(contentsElm, "id", "menu")
	menuBoxElm := findByAttr(menuElm, "class", "menu_box")
	ulElm := findByData(menuBoxElm, "ul")

	var links []string
	for node := ulElm.FirstChild; node != nil; node = node.NextSibling {
		if node.Data != "li" {
			continue
		}

		aElm := findByData(node, "a")
		href := findAttrVal(aElm.Attr, "href")
		if href == "" {
			continue
		}

		links = append(links, href)
	}

	return links
}

func pdfLinks(doc *html.Node) []struct{ title, link string } {
	htmlElm := nthChild(doc, 2)
	bodyElm := findByData(htmlElm, "body")
	pageBodyElm := findByAttr(bodyElm, "id", "pagebody")
	contentsElm := findByAttr(pageBodyElm, "id", "contents")
	mainElm := findByAttr(contentsElm, "id", "main")

	var links []struct{ title, link string }
	for node := mainElm.FirstChild; node != nil; node = node.NextSibling {
		if findAttrVal(node.Attr, "class") != "program_box" {
			continue
		}

		titleElm := findByAttr(node, "class", "title")
		aElm := findByData(titleElm, "a")
		if aElm == nil {
			continue
		}

		title := titleElm.FirstChild.FirstChild.Data
		link := findAttrVal(titleElm.FirstChild.Attr, "href")
		links = append(links, struct{ title, link string }{title, link})
	}

	return links
}

func findByData(n *html.Node, data string) *html.Node {
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == data {
			return i
		}
	}

	return nil
}

func findByAttr(n *html.Node, key, val string) *html.Node {
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		attrVal := findAttrVal(i.Attr, key)
		if attrVal == val {
			return i
		}
	}

	return nil
}

func nthChild(node *html.Node, n int) *html.Node {
	return nthSibling(node.FirstChild, n-1)
}

func nthSibling(node *html.Node, n int) *html.Node {
	sibling := node
	for i := 0; i < n; i++ {
		sibling = sibling.NextSibling
	}
	return sibling
}

func findAttrVal(a []html.Attribute, key string) string {
	for _, a := range a {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}
