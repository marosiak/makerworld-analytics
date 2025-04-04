package domain

import (
	"encoding/json"
	"log"
	"makerworld-analytics/makerworld"
	"time"
)

type Statistics struct {
	TotalPoints      float32 `json:"totalPoints"`
	PointsFromBoosts float32 `json:"pointsFromBoosts"`
	PointsFromDesign float32 `json:"pointsFromDesign"`
	PointsOther      float32 `json:"pointsOther"`

	PointsPerDate PointsPerDateMap `json:"pointsPerDay"`
}

type PointsPerDateMap map[time.Time]float32

func (s PointsPerDateMap) SumPointsChange() float32 {
	var total float32
	for _, points := range s {
		total += points
	}
	return total
}

func (s PointsPerDateMap) AveragePointsChange() float32 {
	var total float32
	var count int
	for _, points := range s {
		total += points
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float32(count)
}

func (s PointsPerDateMap) FilterByDate(start, end *time.Time) PointsPerDateMap {
	filtered := make(PointsPerDateMap)
	for date, points := range s {
		if start != nil && date.Before(*start) {
			continue
		}

		if end != nil && date.After(*end) {
			continue
		}

		filtered[date] = points
	}
	return filtered
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

	pointsPerDay := make(map[time.Time]float32)
	for _, hit := range rawStatsData.Hits {
		if !hit.CreateTime.IsZero() {
			pointsPerDay[hit.CreateTime] += hit.PointChange
		}
	}

	return &Statistics{
		TotalPoints:      rawStatsData.TotalIncome,
		PointsFromBoosts: incomeFromBoostsAllTime,
		PointsFromDesign: incomeFromDesignAllTime,
		PointsOther:      incomeFromInstanceRewardAllTime,
		PointsPerDate:    pointsPerDay,
	}
}
