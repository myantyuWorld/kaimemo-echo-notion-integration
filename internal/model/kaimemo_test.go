package model

import (
	"reflect"
	"testing"
)

func TestKaimemoAmountRecords_GroupByWeek(t *testing.T) {
	tests := []struct {
		name     string
		records  KaimemoAmountRecords
		expected []WeeklySummary
	}{
		{
			name: "empty records",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{},
			},
			expected: []WeeklySummary{},
		},
		{
			name: "single week records",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-15", Amount: 1000},
					{Date: "2023-05-16", Amount: 2000},
					{Date: "2023-05-17", Amount: 3000},
				},
			},
			expected: []WeeklySummary{
				{
					WeekStart:   "2023-05-14",
					WeekEnd:     "2023-05-20",
					TotalAmount: 6000,
					Items: []KaimemoAmount{
						{Date: "2023-05-15", Amount: 1000},
						{Date: "2023-05-16", Amount: 2000},
						{Date: "2023-05-17", Amount: 3000},
					},
				},
			},
		},
		{
			name: "multiple weeks across month boundary",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-30", Amount: 1000},
					{Date: "2023-06-01", Amount: 2000},
					{Date: "2023-06-05", Amount: 3000},
				},
			},
			expected: []WeeklySummary{
				{
					WeekStart:   "2023-05-28",
					WeekEnd:     "2023-06-03",
					TotalAmount: 3000,
					Items: []KaimemoAmount{
						{Date: "2023-05-30", Amount: 1000},
						{Date: "2023-06-01", Amount: 2000},
					},
				},
				{
					WeekStart:   "2023-06-04",
					WeekEnd:     "2023-06-10",
					TotalAmount: 3000,
					Items: []KaimemoAmount{
						{Date: "2023-06-05", Amount: 3000},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.records.GroupByWeek()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\nGroupByWeek() =\n%+v\nwant\n%+v", result, tt.expected)
			}
		})
	}
}
