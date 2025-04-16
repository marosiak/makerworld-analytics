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

type PieDataItem struct {
	Value float32
	Name  string
}
type PieData []PieDataItem

func (p PieData) ToJS() interface{} {
	var arr []interface{}
	for _, v := range p {
		arr = append(arr, map[string]interface{}{"value": v.Value, "name": v.Name})
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

type DataZoom struct {
	Type  string
	Start int
	End   int
}
type ChartOption struct {
	Color    []string
	Legend   LegendOption
	Series   []SeriesOption
	Title    TitleOption
	Toolbox  ToolboxOption
	Tooltip  TooltipOption
	XAxis    []XAxisOption
	YAxis    []YAxisOption
	DataZoom []DataZoom
}

type LegendOption struct {
	Data  []string
	Other map[string]interface{}
}

type ItemStyle struct {
	BorderRadius int
}
type StackType string

const (
	StackTypeTotal StackType = "Total"
)

type Emphasis struct {
	Focus string
}

type SeriesOption struct {
	Name      string
	Type      string
	Data      EChartsData
	Label     map[string]interface{}
	Radius    []string
	PadAngle  int
	ItemStyle ItemStyle
	Stack     StackType
	Emphasis  *Emphasis
	AreaStyle map[string]interface{}
}

type TitleOption struct {
	Text string
	More map[string]interface{}
}

type ToolboxOption struct {
	Show bool
}

type TooltipOption struct {
	Show  bool
	Other map[string]interface{}
}

type XAxisOption struct {
	Data  []string
	Other map[string]interface{}
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
		series := map[string]interface{}{
			"name": s.Name,
			"type": s.Type,
		}
		if s.Data != nil {
			series["data"] = s.Data.ToJS()
		} else {
			series["data"] = []interface{}{}
		}
		if s.Label != nil {
			series["label"] = s.Label
		}
		if len(s.Radius) > 0 {
			series["radius"] = sliceOfStringToInterface(s.Radius)
		}
		if s.PadAngle != 0 {
			series["padAngle"] = s.PadAngle
		}
		if s.Stack != "" {
			series["stack"] = string(s.Stack)
		}
		if s.Emphasis != nil {
			series["emphasis"] = map[string]interface{}{}
			if s.Emphasis.Focus != "" {
				series["emphasis"] = map[string]interface{}{
					"focus": s.Emphasis.Focus,
				}
			}
		}
		if s.AreaStyle != nil {
			series["areaStyle"] = s.AreaStyle
		}
		if s.ItemStyle.BorderRadius != 0 {
			series["itemStyle"] = map[string]interface{}{
				"borderRadius": s.ItemStyle.BorderRadius,
			}
		}
		seriesArr = append(seriesArr, series)
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

	tooltip := make(map[string]interface{})
	if o.Tooltip.Show {
		tooltip["show"] = true
	}

	if o.Tooltip.Other != nil {
		for k, v := range o.Tooltip.Other {
			tooltip[k] = v
		}
	}
	m["tooltip"] = tooltip

	var xAxisArr []interface{}
	for _, xa := range o.XAxis {
		xmap := map[string]interface{}{}
		if len(xa.Data) > 0 {
			xmap["data"] = sliceOfStringToInterface(xa.Data)
		}
		if xa.Other != nil {
			for k, v := range xa.Other {
				xmap[k] = v
			}
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

	if len(o.DataZoom) > 0 {
		dataZoomArr := make([]interface{}, 0, len(o.DataZoom))
		for _, dz := range o.DataZoom {
			dzMap := map[string]interface{}{
				"type":  dz.Type,
				"start": dz.Start,
				"end":   dz.End,
			}
			dataZoomArr = append(dataZoomArr, dzMap)
		}
		m["dataZoom"] = dataZoomArr
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
