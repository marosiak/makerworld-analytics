package domain

import (
	"encoding/json"
	"log"
	"makerworld-analytics/makerworld"
)

type Statistics struct {
	TotalPoints      float32 `json:"totalPoints"`
	PointsFromBoosts float32 `json:"pointsFromBoosts"`
	PointsFromDesign float32 `json:"pointsFromDesign"`
	PointsOther      float32 `json:"pointsOther"`

	PointsPerDay map[string]float32 `json:"pointsPerDay"`
}

type MoneyMultiplier float32

const (
	BankPayoutMultiplier MoneyMultiplier = 0.066
	VouchersMultiplier   MoneyMultiplier = 0.07633587786
)

func (s Statistics) ToEuro(multiplier MoneyMultiplier, pointsAmount float32) float32 {
	return pointsAmount * float32(multiplier)
}

func NewStatistics(sourceJson string) *Statistics {
	rawStatsData := makerworld.PointsStatistics{}
	err := json.Unmarshal([]byte(sourceJson), &rawStatsData)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s\n", err)
		return nil
	}

	incomeFromBoostsAllTime := rawStatsData.Hits.FilterByType(makerworld.RevenueSourceBoost).SumPointsChange()
	incomeFromDesignAllTime := rawStatsData.Hits.FilterByType(makerworld.RevenueSourceDesignReward).SumPointsChange()
	incomeFromInstanceRewardAllTime := rawStatsData.Hits.FilterByType(makerworld.RevenueSourceInstanceReward).SumPointsChange()

	pointsPerDay := make(map[string]float32)
	for _, hit := range rawStatsData.Hits {
		if !hit.CreateTime.IsZero() {
			pointsPerDay[hit.CreateTime.Format("2006-01-02")] += hit.PointChange
		}
	}

	return &Statistics{
		TotalPoints:      rawStatsData.TotalIncome,
		PointsFromBoosts: incomeFromBoostsAllTime,
		PointsFromDesign: incomeFromDesignAllTime,
		PointsOther:      incomeFromInstanceRewardAllTime,
		PointsPerDay:     pointsPerDay,
	}
}
