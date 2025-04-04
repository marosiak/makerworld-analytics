package echarts_wasm

type EChartsData interface {
	ToJS() interface{}
}

type BarData struct {
	Values []float32
}

func (b BarData) ToJS() interface{} {
	var arr []interface{}
	for _, v := range b.Values {
		arr = append(arr, map[string]interface{}{"value": v})
	}
	return arr
}

type NumericData struct {
	Values []float32
}

func (n NumericData) ToJS() interface{} {
	var arr []interface{}
	for _, v := range n.Values {
		arr = append(arr, v)
	}
	return arr
}

type StringData struct {
	Values []string
}

func (s StringData) ToJS() interface{} {
	var arr []interface{}
	for _, v := range s.Values {
		arr = append(arr, v)
	}
	return arr
}

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
	Data  EChartsData
	Label map[string]interface{}
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
		sm := map[string]interface{}{
			"name": s.Name,
			"type": s.Type,
		}
		if s.Data != nil {
			sm["data"] = s.Data.ToJS()
		} else {
			sm["data"] = []interface{}{}
		}
		if s.Label != nil {
			sm["label"] = s.Label
		}
		seriesArr = append(seriesArr, sm)
	}
	if len(seriesArr) > 0 {
		m["series"] = seriesArr
	} else {
		delete(m, "series")
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
