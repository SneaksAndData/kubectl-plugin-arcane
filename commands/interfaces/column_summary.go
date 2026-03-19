package interfaces

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DowntimeSummary defines an interface for summarizing downtime information, including counts and details of downtime events.
type DowntimeSummary interface {

	// Counts returns a table summarizing the counts of downtime events, categorized by relevant criteria.
	Counts() *v1.Table

	// CountsRaw returns a raw map of counts of downtime events, categorized by relevant criteria, without any formatting.
	CountsRaw() map[string]int

	// Details returns a table containing detailed information about individual downtime events, including relevant metadata and annotations.
	Details() *v1.Table
}
