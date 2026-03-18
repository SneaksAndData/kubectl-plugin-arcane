package interfaces

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DowntimeSummary interface {
	Counts() *v1.Table
}
