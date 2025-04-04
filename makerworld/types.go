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
)

type Hit struct {
	Type                      RevenueSource `json:"type"`
	PointChange               float32       `json:"pointChange"`
	PointChangeRegular        float32       `json:"pointChangeRegular"`
	PointChangeExclusive      float32       `json:"pointChangeExclusive"`
	IsExclusiveBonus          bool          `json:"isExclusiveBonus"`
	PointTime                 time.Time     `json:"pointTime"`
	CreateTime                time.Time     `json:"createTime"`
	SyncSource                string        `json:"syncSource"`
	ExtInfoBoostExchangePoint struct {
		DesignId       int    `json:"designId"`
		DesignTitle    string `json:"designTitle"`
		FromUid        int64  `json:"fromUid"`
		FromUsername   string `json:"fromUsername"`
		FromUserHandle string `json:"fromUserHandle"`
	} `json:"extInfoBoostExchangePoint,omitempty"`
	ExtInfoMakerLabFirstExport         interface{} `json:"extInfoMakerLabFirstExport"`
	ExtInfoAcademyCourseAward          interface{} `json:"extInfoAcademyCourseAward"`
	ExtinfoOperationCancelExchangeCash interface{} `json:"extinfoOperationCancelExchangeCash"`
	InstanceReward                     struct {
		DesignId              int    `json:"designId"`
		DesignTitle           string `json:"designTitle"`
		InstanceId            int    `json:"instanceId"`
		InstanceTitle         string `json:"instanceTitle"`
		DownloadAndPrintCount int    `json:"downloadAndPrintCount"`
		LikeAndRateCount      int    `json:"likeAndRateCount"`
		AverageRatingGe       string `json:"averageRatingGe"`
		DonatedByUid          int    `json:"donatedByUid"`
		DonatedByUsername     string `json:"donatedByUsername"`
	} `json:"instanceReward,omitempty"`
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
