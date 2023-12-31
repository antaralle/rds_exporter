package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type mCollector struct {
	rds_metric *prometheus.Desc
}

func newCollector() *mCollector {

	return &mCollector{
		rds_metric: prometheus.NewDesc("rds_metric",
			"info about rds connections",
			[]string{"ClientAddress", "ConnectedResource", "UserName", "ProtocolName", "TunnelId"}, nil,
		),
	}
}

func (collector *mCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.rds_metric

}

func (collector *mCollector) Collect(ch chan<- prometheus.Metric) {

	var metricValue = 1.0
	var labels []string
	var n []string

	res := []rdsUserStats{}
	wmi.QueryNamespace("SELECT ClientAddress, ConnectedResource, UserName, ProtocolName, TunnelId  FROM Win32_TSGatewayConnection", &res, "root\\cimv2\\TerminalServices")

	for _, v := range res {
		labels = append(n, v.ClientAddress, v.ConnectedResource, v.UserName, v.ProtocolName, v.TunnelId)
		ch <- prometheus.MustNewConstMetric(collector.rds_metric, prometheus.GaugeValue, metricValue, labels...)
	}

}

type rdsUserStats struct {
	ClientAddress     string
	ConnectedResource string
	UserName          string
	ProtocolName      string
	TunnelId          string
}

type RdsStats struct {
	Users            []rdsUserStats
	totalConnections int
}

func main() {
	c := newCollector()
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
