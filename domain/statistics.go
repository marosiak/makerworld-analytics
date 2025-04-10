package domain

import (
	"encoding/json"
	"log"
	"makerworld-analytics/makerworld"
	"time"
)

type PointAssignment struct {
	PointChange float32   `json:"pointChange"`
	CreateTime  time.Time `json:"createTime"`
}

type PointsAssignmentList []PointAssignment

func (p PointsAssignmentList) SumPointsChange() float32 {
	var total float32
	for _, points := range p {
		total += points.PointChange
	}
	return total
}

func (p PointsAssignmentList) FilterDate(start, end *time.Time) PointsAssignmentList {
	filtered := make(PointsAssignmentList, 0)
	for _, points := range p {
		if start != nil && points.CreateTime.Before(*start) {
			continue
		}

		if end != nil && points.CreateTime.After(*end) {
			continue
		}

		filtered = append(filtered, points)
	}
	return filtered
}

type PointsPerDesign map[DesignID]PointsAssignmentList
type Statistics struct {
	TotalPoints      float32 `json:"totalPoints"`
	PointsFromBoosts float32 `json:"pointsFromBoosts"`
	PointsFromDesign float32 `json:"pointsFromDesign"`
	PointsOther      float32 `json:"pointsOther"`

	PointsPerDate   PointsPerDateMap `json:"pointsPerDay"`
	PointsPerDesign PointsPerDesign  `json:"pointsPerModel"`

	AllPublishedDesigns []PublishedDesign `json:"allPublishedDesigns"`
}

func (s Statistics) GetDesignByID(id DesignID) (PublishedDesign, bool) {
	for _, design := range s.AllPublishedDesigns {
		if design.ID == id {
			return design, true
		}
	}
	return PublishedDesign{}, false
}

type PointsPerModel map[string]float32
type PointsPerDateMap map[time.Time]float32

func (s PointsPerDateMap) SumPointsChange() float32 {
	var total float32
	for _, points := range s {
		total += points
	}
	return total
}

func (s PointsPerDateMap) AveragePointsChange() float32 {
	if len(s) == 0 {
		return 0
	}

	normalized := map[time.Time]float32{}
	firstDay := time.Time{}
	lastDay := time.Time{}

	for date, points := range s {
		d := date.Format("2006-01-02")
		dateOnly, _ := time.Parse("2006-01-02", d)

		normalized[dateOnly] += points

		if firstDay.IsZero() || dateOnly.Before(firstDay) {
			firstDay = dateOnly
		}
		if lastDay.IsZero() || dateOnly.After(lastDay) {
			lastDay = dateOnly
		}
	}

	days := int(lastDay.Sub(firstDay).Hours()/24) + 1
	if days <= 1 {
		return 0
	}

	pointsPerDay := make([]float32, days)
	for i := 0; i < days; i++ {
		day := firstDay.AddDate(0, 0, i)
		pointsPerDay[i] = normalized[day]
	}

	totalChange := float32(0)
	for i := 1; i < days; i++ {
		change := pointsPerDay[i] - pointsPerDay[i-1]
		if change < 0 {
			change = -change
		}
		totalChange += change
	}

	averageChange := totalChange / float32(days-1)
	return averageChange
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

	pointsPerDayMap := make(map[time.Time]float32)
	pointsPerDesignMap := make(PointsPerDesign)
	allPublishedDesigns := make(PublishedDesignsList, 0)

	for _, hit := range rawStatsData.Hits {
		if !hit.CreateTime.IsZero() {
			pointsPerDayMap[hit.CreateTime] += hit.PointChange

			pointsAssignment := PointAssignment{
				PointChange: hit.PointChange,
				CreateTime:  hit.CreateTime,
			}

			designId := DesignID(hit.DesignId())

			_, exists := pointsPerDesignMap[designId]
			if exists {
				pointsPerDesignMap[designId] = append(pointsPerDesignMap[designId], pointsAssignment)
			} else {
				pointsPerDesignMap[designId] = PointsAssignmentList{
					pointsAssignment,
				}
			}

			if !allPublishedDesigns.Exists(designId) {
				allPublishedDesigns = append(allPublishedDesigns, PublishedDesign{
					ID:   designId,
					Name: hit.DesignName(),
				})
			}

		}
	}

	return &Statistics{
		TotalPoints:         rawStatsData.TotalIncome,
		PointsFromBoosts:    incomeFromBoostsAllTime,
		PointsFromDesign:    incomeFromDesignAllTime,
		PointsOther:         incomeFromInstanceRewardAllTime,
		PointsPerDate:       pointsPerDayMap,
		PointsPerDesign:     pointsPerDesignMap,
		AllPublishedDesigns: allPublishedDesigns,
	}
}
