/*
Obtains Google Analytics RealTime API metrics, and presents them to
prometheus for scraping.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/analytics/v3"
)

var (
	promGauge = make(map[string]prometheus.Gauge)
	Metrics   = []string{}
	reg       = new(prometheus.Registry)
)

func init() {
	loadMetrics()
	reg = prometheus.NewRegistry()

	// All metrics are registered as Prometheus Gauge
	for _, metric := range Metrics {
		promGauge[metric] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        fmt.Sprintf("ga_%s", strings.Replace(metric, ":", "_", 1)),
			Help:        fmt.Sprintf("Google Analytics %s", metric),
			ConstLabels: map[string]string{"job": "googleAnalytics"},
		})

		reg.MustRegister(promGauge[metric])
	}
}

func main() {
	creds := getCreds()

	// JSON web token configuration
	jwtc := jwt.Config{
		Email:        creds["client_email"],
		PrivateKey:   []byte(creds["private_key"]),
		PrivateKeyID: creds["private_key_id"],
		Scopes:       []string{analytics.AnalyticsReadonlyScope},
		TokenURL:     creds["token_uri"],
		// Expires:      time.Duration(1) * time.Hour, // Expire in 1 hour
	}

	httpClient := jwtc.Client(oauth2.NoContext)
	as, err := analytics.New(httpClient)
	if err != nil {
		panic(err)
	}

	// Authenticated RealTime Google Analytics API service
	rts := analytics.NewDataRealtimeService(as)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry: reg,
	}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Google Analytics Exporter</title></head>
			<body>
			<h1>Google Analytics Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	go http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("GA_PORT")), nil)

	for {
		for _, metric := range Metrics {
			// Go routine per mertic
			go func(metric string) {
				val := getMetric(rts, metric)
				// Gauge value to float64
				valf, _ := strconv.ParseFloat(val, 64)
				promGauge[metric].Set(valf)
			}(metric)
		}

		interval, err := strconv.ParseInt(os.Getenv("GA_INTERVAL"), 10, 16)
		if err != nil {
			interval = 60
		}

		time.Sleep(time.Second * time.Duration(interval))
	}
}

// getMetric queries GA RealTime API for a specific metric.
func getMetric(rts *analytics.DataRealtimeService, metric string) string {
	getc := rts.Get(os.Getenv("GA_VIEWID"), metric)
	m, err := getc.Do()
	if err != nil {
		panic(err)
	}

	if len(m.Rows) == 0 {
		return ""
	}

	if len(m.Rows[0]) == 0 {
		return ""
	}

	return m.Rows[0][0]
}

func loadMetrics() {
	raw := os.Getenv("GA_METRICS")
	Metrics = strings.Fields(raw)
}

// https://console.developers.google.com/apis/credentials
// 'Service account keys' creds formated file is expected.
// NOTE: the email from the creds has to be added to the Analytics permissions
func getCreds() (r map[string]string) {
	data, err := ioutil.ReadFile("/run/secrets/GA_CREDS")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	return r
}
