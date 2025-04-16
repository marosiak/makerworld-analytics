package domain

import (
	"fmt"
	"time"
)

type Settings struct {
	StartDate         *time.Time
	EndDate           *time.Time
	MoneyMultiplier   MoneyMultiplier
	PublicationFilter *PublishedDesign
}

func (s Settings) String() string {
	return fmt.Sprintf("[Settings] StartDate: %s, EndDate: %s, MoneyMultiplier: %s", s.StartDate, s.EndDate, s.MoneyMultiplier)
}
