package echarts_wasm

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ChartOption struct {
	Color   []string
	Legend  LegendOption
	Series  []SeriesOption
	Title   TitleOption
	Toolbox ToolboxOption
	Tooltip TooltipOption
	XAxis   []XAxisOption
	YAxis   []YAxisOption
}

type LegendOption struct {
	Data  []string
	Other map[string]interface{}
}

type SeriesOption struct {
	Name  string
	Type  string
	Data  []DataItem
	Label map[string]interface{}
}

type DataItem struct {
	Value interface{}
}

type TitleOption struct {
	Text string
	More map[string]interface{}
}

type ToolboxOption struct {
	Show bool
}

type TooltipOption struct {
	Show bool
}

type XAxisOption struct {
	Data []string
}

type YAxisOption struct {
	Some map[string]interface{}
}

func (o ChartOption) ToValue() app.Value {
	m := map[string]interface{}{}
	if len(o.Color) > 0 {
		m["color"] = o.Color
	}
	l := map[string]interface{}{}
	if len(o.Legend.Data) > 0 {
		l["data"] = o.Legend.Data
	}
	if o.Legend.Other != nil {
		for k, v := range o.Legend.Other {
			l[k] = v
		}
	}
	m["legend"] = l
	var seriesArr []interface{}
	for _, s := range o.Series {
		sm := map[string]interface{}{}
		sm["name"] = s.Name
		sm["type"] = s.Type
		var dataArr []interface{}
		for _, d := range s.Data {
			dataArr = append(dataArr, map[string]interface{}{"value": d.Value})
		}
		sm["data"] = dataArr
		if s.Label != nil {
			sm["label"] = s.Label
		}
		seriesArr = append(seriesArr, sm)
	}
	m["series"] = seriesArr
	t := map[string]interface{}{}
	if o.Title.Text != "" {
		t["text"] = o.Title.Text
	}
	if o.Title.More != nil {
		for k, v := range o.Title.More {
			t[k] = v
		}
	}
	m["title"] = t
	m["toolbox"] = map[string]interface{}{
		"show": o.Toolbox.Show,
	}
	m["tooltip"] = map[string]interface{}{
		"show": o.Tooltip.Show,
	}
	var xAxisArr []interface{}
	for _, x := range o.XAxis {
		xmap := map[string]interface{}{}
		if len(x.Data) > 0 {
			xmap["data"] = x.Data
		}
		xAxisArr = append(xAxisArr, xmap)
	}
	m["xAxis"] = xAxisArr
	var yAxisArr []interface{}
	for _, y := range o.YAxis {
		ymap := map[string]interface{}{}
		if y.Some != nil {
			for k, v := range y.Some {
				ymap[k] = v
			}
		}
		yAxisArr = append(yAxisArr, ymap)
	}
	m["yAxis"] = yAxisArr
	return app.ValueOf(m)
}

func (o ChartOption) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"color":   sliceOfStringToInterface(o.Color),
		"legend":  map[string]interface{}{},
		"series":  []interface{}{},
		"title":   map[string]interface{}{},
		"toolbox": map[string]interface{}{},
		"tooltip": map[string]interface{}{},
		"xAxis":   []interface{}{},
		"yAxis":   []interface{}{},
	}

	if len(o.Color) == 0 {
		delete(m, "color")
	}

	leg := map[string]interface{}{}
	if len(o.Legend.Data) > 0 {
		leg["data"] = sliceOfStringToInterface(o.Legend.Data)
	}
	if o.Legend.Other != nil {
		for k, v := range o.Legend.Other {
			leg[k] = v
		}
	}
	if len(leg) > 0 {
		m["legend"] = leg
	} else {
		delete(m, "legend")
	}

	var seriesArr []interface{}
	for _, s := range o.Series {
		seriesMap := map[string]interface{}{
			"name": s.Name,
			"type": s.Type,
			"data": []interface{}{},
		}
		if len(s.Data) > 0 {
			var d []interface{}
			for _, dp := range s.Data {
				d = append(d, map[string]interface{}{"value": dp.Value})
			}
			seriesMap["data"] = d
		}
		if s.Label != nil {
			seriesMap["label"] = s.Label
		}
		seriesArr = append(seriesArr, seriesMap)
	}
	if len(seriesArr) > 0 {
		m["series"] = seriesArr
	}

	titleMap := map[string]interface{}{}
	if o.Title.Text != "" {
		titleMap["text"] = o.Title.Text
	}
	if o.Title.More != nil {
		for k, v := range o.Title.More {
			titleMap[k] = v
		}
	}
	if len(titleMap) > 0 {
		m["title"] = titleMap
	} else {
		delete(m, "title")
	}

	if o.Toolbox.Show {
		m["toolbox"] = map[string]interface{}{"show": true}
	} else {
		delete(m, "toolbox")
	}

	if o.Tooltip.Show {
		m["tooltip"] = map[string]interface{}{"show": true}
	} else {
		delete(m, "tooltip")
	}

	var xAxisArr []interface{}
	for _, xa := range o.XAxis {
		xmap := map[string]interface{}{}
		if len(xa.Data) > 0 {
			xmap["data"] = sliceOfStringToInterface(xa.Data)
		}
		xAxisArr = append(xAxisArr, xmap)
	}
	if len(xAxisArr) > 0 {
		m["xAxis"] = xAxisArr
	} else {
		delete(m, "xAxis")
	}

	var yAxisArr []interface{}
	for _, ya := range o.YAxis {
		ymap := map[string]interface{}{}
		if ya.Some != nil {
			for k, v := range ya.Some {
				ymap[k] = v
			}
		}
		yAxisArr = append(yAxisArr, ymap)
	}
	if len(yAxisArr) > 0 {
		m["yAxis"] = yAxisArr
	} else {
		delete(m, "yAxis")
	}

	return m
}

func sliceOfStringToInterface(ss []string) []interface{} {
	arr := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		arr = append(arr, s)
	}
	return arr
}
