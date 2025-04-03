package go_echarts_bridge

import (
	"bytes"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/net/html"
	"io"
)

func ComponentFromChart(chartRenderer render.Renderer) app.UI {
	buff := bytes.NewBuffer(nil)
	err := chartRenderer.Render(buff)
	if err != nil {
		return app.Div().Text(fmt.Sprintf("Error rendering chart: %v", err))
	}

	rawHtml, err := io.ReadAll(buff)
	if err != nil {
		return app.Div().Text(fmt.Sprintf("Error reading HTML: %v", err))
	}

	doc, err := html.Parse(bytes.NewReader(rawHtml))
	if err != nil {
		return app.Div().Text(fmt.Sprintf("Error parsing HTML: %v", err))
	}

	var bodyNode *html.Node

	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			bodyNode = n
			return
		}
		for c := n.FirstChild; c != nil && bodyNode == nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)

	if bodyNode == nil {
		return app.Div().Text("No body found in HTML")
	}

	var buf bytes.Buffer
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return app.Raw(buf.String())
	//return app.Div().ID("chart-component").Text(
	//	buf.String())
}
