package makerworld

import (
	"sort"
	"time"
)

type RevenueSource string

const (
	RevenueSourceBoost          = "boost_exchange_point"
	RevenueSourceInstanceReward = "instance_reward"
	RevenueSourceDesignReward   = "design_reward"
	RevenueSourceCreateInstance = "create_instance"
)

type Hit struct {
	Type                 RevenueSource `json:"type"`
	PointChange          float32       `json:"pointChange"`
	PointChangeRegular   float32       `json:"pointChangeRegular"`
	PointChangeExclusive float32       `json:"pointChangeExclusive"`
	IsExclusiveBonus     bool          `json:"isExclusiveBonus"`
	PointTime            time.Time     `json:"pointTime"`
	CreateTime           time.Time     `json:"createTime"`
	SyncSource           string        `json:"syncSource"`
	DesignReward         struct {
		ID                    int    `json:"id"`
		Title                 string `json:"title"`
		ModelSource           int    `json:"modelSource"`
		DownloadAndPrintCount int    `json:"downloadAndPrintCount"`
		LikeAndRateCount      int    `json:"likeAndRateCount"`
	} `json:"designReward"`
	ExtInfoBoostExchangePoint struct {
		DesignID       int    `json:"designId"`
		DesignTitle    string `json:"designTitle"`
		FromUID        int64  `json:"fromUid"`
		FromUsername   string `json:"fromUsername"`
		FromUserHandle string `json:"fromUserHandle"`
	} `json:"extInfoBoostExchangePoint,omitempty"`
	ExtInfoMakerLabFirstExport         interface{} `json:"extInfoMakerLabFirstExport"`
	ExtInfoAcademyCourseAward          interface{} `json:"extInfoAcademyCourseAward"`
	ExtinfoOperationCancelExchangeCash interface{} `json:"extinfoOperationCancelExchangeCash"`
	InstanceReward                     struct {
		DesignID              int    `json:"designId"`
		DesignTitle           string `json:"designTitle"`
		InstanceID            int    `json:"instanceId"`
		InstanceTitle         string `json:"instanceTitle"`
		DownloadAndPrintCount int    `json:"downloadAndPrintCount"`
		LikeAndRateCount      int    `json:"likeAndRateCount"`
		AverageRatingGe       string `json:"averageRatingGe"`
		DonatedByUID          int    `json:"donatedByUid"`
		DonatedByUsername     string `json:"donatedByUsername"`
	} `json:"instanceReward,omitempty"`

	Instance struct {
		DesignId          int    `json:"designId"`
		DesignTitle       string `json:"designTitle"`
		InstanceId        int    `json:"instanceId"`
		InstanceTitle     string `json:"instanceTitle"`
		DonatedByUid      int    `json:"donatedByUid"`
		DonatedByUsername string `json:"donatedByUsername"`
	} `json:"instance"`
	Design struct {
		Id          int    `json:"id"`
		Title       string `json:"title"`
		ModelSource int    `json:"modelSource"`
	} `json:"design"`
}

func (h Hit) DesignID() int {
	tmp := h.ExtInfoBoostExchangePoint.DesignID
	if tmp != 0 {
		return tmp
	}

	tmp = h.InstanceReward.DesignID
	if tmp != 0 {
		return tmp
	}

	tmp = h.Instance.DesignId
	if tmp != 0 {
		return tmp
	}

	tmp = h.DesignReward.ID
	if tmp != 0 {
		return tmp
	}

	return h.Design.Id
}

func (h Hit) DesignName() string {
	tmp := h.ExtInfoBoostExchangePoint.DesignTitle
	if tmp != "" {
		return tmp
	}

	tmp = h.InstanceReward.DesignTitle
	if tmp != "" {
		return tmp
	}

	tmp = h.Instance.DesignTitle
	if tmp != "" {
		return tmp
	}

	tmp = h.DesignReward.Title
	if tmp != "" {
		return tmp
	}

	tmp = h.Design.Title
	if tmp != "" {
		return tmp
	}

	return "[Unknown source]"
}

type HitsList []Hit
type PointsStatistics struct {
	Total                 float32  `json:"total"`
	TotalIncome           float32  `json:"totalIncome"`
	TotalExpense          float32  `json:"totalExpense"`
	TotalRegularIncome    float32  `json:"totalRegularIncome"`
	TotalRegularExpense   float32  `json:"totalRegularExpense"`
	TotalExclusiveIncome  float64  `json:"totalExclusiveIncome"`
	TotalExclusiveExpense float32  `json:"totalExclusiveExpense"`
	Hits                  HitsList `json:"hits"`
}

func (a HitsList) SortByDate(ascent bool) HitsList {
	sort.Slice(a, func(i, j int) bool {
		if ascent {
			return a[i].PointTime.Before(a[j].PointTime)
		}
		return a[i].PointTime.After(a[j].PointTime)
	})
	return a
}

func (a HitsList) FilterByTimeRange(start time.Time, end time.Time) []Hit {
	var hits []Hit
	for _, hit := range a {
		if hit.PointTime.After(start) && hit.PointTime.Before(end) {
			hits = append(hits, hit)
		}
	}
	return hits
}

func (a HitsList) FilterByType(revenueSource RevenueSource) HitsList {
	var hits []Hit
	for _, hit := range a {
		if hit.Type == revenueSource {
			hits = append(hits, hit)
		}
	}
	return hits
}

func (a HitsList) SumPointsChange() float32 {
	var sum float32
	for _, hit := range a {
		sum += hit.PointChange
	}
	return sum
}
