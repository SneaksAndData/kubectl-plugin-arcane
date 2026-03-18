package services

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ interfaces.DowntimeSummary = (*DowntimeSummary)(nil)

type DowntimeSummary struct {
	counts map[string]int
}

func NewDowntimeSummary(counts map[string]int) *DowntimeSummary {
	return &DowntimeSummary{counts: counts}
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

	for key, count := range d.counts {
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
	return d.counts
}
