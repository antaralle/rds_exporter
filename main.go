package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type mCollector struct {
	rds_metric       *prometheus.Desc
	rds_total_metric *prometheus.Desc
}

func newCollector() *mCollector {

	return &mCollector{
		rds_metric: prometheus.NewDesc("rds_metric",
			"info about rds connections",
			[]string{"ClientAddress", "ConnectedResource", "UserName"}, nil,
		),
		rds_total_metric: prometheus.NewDesc("rds_total_metric",
			"total rds connections",
			nil, nil,
		),
	}
}

func (collector *mCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.rds_metric
	ch <- collector.rds_total_metric
}

func (collector *mCollector) Collect(ch chan<- prometheus.Metric) {

	var metricValue = 1.0
	var labels []string

	res := []rdsUserStats{}
	wmi.QueryNamespace("SELECT ClientAddress, ConnectedResource, UserName FROM Win32_TSGatewayConnection", &res, "root\\cimv2\\TerminalServices")
	labels = append(labels, res[0].ClientAddress)
	labels = append(labels, res[0].ConnectedResource)
	labels = append(labels, res[0].UserName)

	// m1 := prometheus.MustNewConstMetric(collector.rds_metric, prometheus.GaugeValue, metricValue, labels...)
	// m2 := prometheus.MustNewConstMetric(collector.rds_total_metric, prometheus.GaugeValue, float64(len(res)))

	ch <- prometheus.MustNewConstMetric(collector.rds_metric, prometheus.GaugeValue, metricValue, labels...)
	ch <- prometheus.MustNewConstMetric(collector.rds_total_metric, prometheus.GaugeValue, float64(len(res)))
	var l2 []string
	l2 = append(labels, res[1].ClientAddress)
	l2 = append(labels, res[1].ConnectedResource)
	l2 = append(labels, res[1].UserName)
	ch <- prometheus.MustNewConstMetric(collector.rds_metric, prometheus.GaugeValue, metricValue, l2...)

	//ch <- prometheus.MustNewConstMetric(collector.rds_metric, prometheus.GaugeValue, metricValue, labels...)
	ch <- prometheus.MustNewConstMetric(collector.rds_total_metric, prometheus.GaugeValue, float64(len(res)))
}

type rdsUserStats struct {
	ClientAddress     string
	ConnectedResource string
	UserName          string
}

type RdsStats struct {
	Users            []rdsUserStats
	totalConnections int
}

func getRdsStatistics() RdsStats {
	var res []rdsUserStats

	wmi.QueryNamespace("SELECT ClientAddress, ConnectedResource, UserName FROM Win32_TSGatewayConnection", &res, "root\\cimv2\\TerminalServices")

	return RdsStats{
		Users:            res,
		totalConnections: len(res),
	}
}
func main() {
	c := newCollector()
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	r := getRdsStatistics()
	fmt.Println(r)

}
