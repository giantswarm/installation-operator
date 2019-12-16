package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelInstallation = "installation"
)

var (
	ScheduleDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("installation_operator", "installation", "info"),
		"Todo description of the installation operator installation metric",
		[]string{
			labelInstallation,
		},
		nil,
	)
)

type InstallationConfig struct {
}

type Installation struct {
}

func NewInstallation(config InstallationConfig) (*Installation, error) {
	r := &Installation{}

	return r, nil
}

func (r *Installation) Collect(ch chan<- prometheus.Metric) error {
	return nil
}

func (r *Installation) Describe(ch chan<- *prometheus.Desc) error {
	ch <- ScheduleDesc

	return nil
}
