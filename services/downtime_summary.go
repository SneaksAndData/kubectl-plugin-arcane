package services

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ interfaces.DowntimeSummary = (*DowntimeSummary)(nil)

type DowntimeSummary struct {
	groupedByKey map[string][]string
}

func NewDowntimeSummary(counts map[string][]string) *DowntimeSummary {
	return &DowntimeSummary{groupedByKey: counts}
}

func (d *DowntimeSummary) Counts() *metav1.Table { // coverage-ignore (tested in integration tests)
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Table",
			APIVersion: "meta.k8s.io/v1",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Count", Type: "integer"},
		},
	}

	for key, count := range d.groupedByKey {
		row := metav1.TableRow{
			Cells: []interface{}{
				key,
				count,
			},
		}
		table.Rows = append(table.Rows, row)
	}

	return table
}

func (d *DowntimeSummary) CountsRaw() map[string]int {
	counts := make(map[string]int)
	for key, items := range d.groupedByKey {
		counts[key] = len(items)
	}

	return counts
}
