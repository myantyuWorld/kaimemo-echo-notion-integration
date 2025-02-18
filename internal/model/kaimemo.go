package model

import (
	"fmt"
	"sort"
	"time"
)

type KaimemoResponse struct {
	ID   string `json:"id"`
	Tag  string `json:"tag"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateKaimemoRequest struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

type KaimemoAmountResponse struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CreateKaimemoAmountRequest struct {
	Date   string `json:"date"`
	Tag    string `json:"tag"`
	Amount int    `json:"amount"`
}

type KaimemoAmount struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Tag    string `json:"tag"`
	Amount int    `json:"amount"`
}

type KaimemoAmountRecords struct {
	Records []KaimemoAmount
}

type WeeklySummary struct {
	WeekStart   string          `json:"weekStart"`
	WeekEnd     string          `json:"weekEnd"`
	TotalAmount int             `json:"totalAmount"`
	Items       []KaimemoAmount `json:"items"`
}

func (k KaimemoAmountRecords) GroupByWeek() []WeeklySummary {
	summaries := make(map[string]*WeeklySummary)

	for _, amount := range k.Records {
		date, _ := time.Parse("2006-01-02", amount.Date)
		year, week := date.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%d", year, week)

		if _, exists := summaries[weekKey]; !exists {
			// 週の開始日を計算
			weekStart := date.AddDate(0, 0, -int(date.Weekday()))
			weekEnd := weekStart.AddDate(0, 0, 6)

			summaries[weekKey] = &WeeklySummary{
				WeekStart:   weekStart.Format("2006-01-02"),
				WeekEnd:     weekEnd.Format("2006-01-02"),
				TotalAmount: 0,
				Items:       []KaimemoAmount{},
			}
		}

		summary := summaries[weekKey]
		summary.TotalAmount += amount.Amount
		summary.Items = append(summary.Items, amount)
	}

	// マップをスライスに変換
	result := make([]WeeklySummary, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, *summary)
	}

	// 日付順にソート
	sort.Slice(result, func(i, j int) bool {
		return result[i].WeekStart < result[j].WeekStart
	})

	return result
}
