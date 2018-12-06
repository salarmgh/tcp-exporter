package main

import (
	"net/http"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"time"
	"os"
	"regexp"
)



type webCollector struct {
	upMetric *prometheus.Desc
	portMetric *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newWebCollector() *webCollector {
	return &webCollector{
		upMetric: prometheus.NewDesc("port_is_up",
			"Shows app is up",
			nil, nil,
		),
		portMetric: prometheus.NewDesc("health_status",
			"Shows the port is up",
			nil, nil,
		),
	}
}

func (collector *webCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.upMetric
	ch <- collector.portMetric
}

func (collector *webCollector) Collect(ch chan<- prometheus.Metric) {

	var serverAddr string
	var status int
	var upValue float64

	serverAddr = os.Getenv("HOST_ADDR")

	rp := regexp.MustCompile("(http.?://)+")
	serverIP := rp.ReplaceAllString(serverAddr, "")

	resp, err := http.Head(serverAddr + "/v1/utilities/health")
	if err != nil {
		status = 0
	} else {
		status = resp.StatusCode
		defer resp.Body.Close()
	}

	if status == 200 {
		upValue = 1
	} else {
		upValue = 0
	}

	var portValue float64
	conn, err := net.DialTimeout("tcp", serverIP, 6 * time.Second)
	if err != nil {
		fmt.Println(err)	
	}
	if conn != nil {
		conn.Close()
		portValue = 1
	} else {
		portValue = 0
	}

	ch <- prometheus.MustNewConstMetric(collector.upMetric, prometheus.CounterValue, upValue)
	ch <- prometheus.MustNewConstMetric(collector.portMetric, prometheus.CounterValue, portValue)
}

func main() {
	web := newWebCollector()
	prometheus.MustRegister(web)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Beginning to serve on port :8283")
	fmt.Println(http.ListenAndServe(":8283", nil))
}
