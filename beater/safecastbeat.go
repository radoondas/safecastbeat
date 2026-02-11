package beater

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	cfglib "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"

	"github.com/radoondas/safecastbeat/config"
)

// Safecastbeat configuration.
type Safecastbeat struct {
	done       chan struct{}
	config     config.Config
	client     beat.Client
	httpClient *http.Client
	logger     *logp.Logger
}

type SafecastResponseData struct {
	Id                    int     `json:"id"`
	User_id               int     `json:"user_id"`
	Value                 float32 `json:"value"`
	Unit                  string  `json:"unit"`
	Location_name         string  `json:"location_name"`
	Device_id             int     `json:"device_id"`
	Original_id           int     `json:"original_id"`
	Measurement_import_id int     `json:"measurement_import_id"`
	Captured_at           string  `json:"captured_at"`
	Height                int     `json:"height"`
	Devicetype_id         string  `json:"devicetype_id"`
	Sensor_id             int     `json:"sensor_id"`
	Station_id            int     `json:"station_id"`
	Channel_id            int     `json:"channel_id"`
	Latitude              float32 `json:"latitude"`
	Longitude             float32 `json:"longitude"`
}

const (
	safecastUrl  = "https://api.safecast.org"
	safecastPath = "/measurements.json"
	selector     = "safecast"

	dateTimeFormat = "2006-01-02+15:04:05"
	perPage        = "1000"
)

// New creates an instance of safecastbeat.
func New(b *beat.Beat, cfg *cfglib.C) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	httpClient, err := buildHTTPClient(c.SSL.CAfile)
	if err != nil {
		return nil, fmt.Errorf("building HTTP client: %w", err)
	}

	logger := logp.NewLogger(selector)

	if c.SSL.CAfile == "" {
		logger.Info("CAFile IS NOT set.")
	}

	bt := &Safecastbeat{
		done:       make(chan struct{}),
		config:     c,
		httpClient: httpClient,
		logger:     logger,
	}

	return bt, nil
}

// buildHTTPClient creates a reusable HTTP client with optional CA certificate.
func buildHTTPClient(caFile string) (*http.Client, error) {
	tlsConfig := &tls.Config{}

	if caFile != "" {
		pemData, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("reading CA file %s: %w", caFile, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pemData) {
			return nil, fmt.Errorf("parsing PEM data from %s", caFile)
		}
		tlsConfig.RootCAs = pool
	}

	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
		Timeout:   30 * time.Second,
	}, nil
}

// Run starts safecastbeat.
func (bt *Safecastbeat) Run(b *beat.Beat) error {
	bt.logger.Info("safecastbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	bt.logger.Infof("Configuration: %+v", bt.config)

	go func() {
		ticker := time.NewTicker(bt.config.Period)
		defer ticker.Stop()

		// UTC is required
		lastTime := time.Now().UTC().Add(-bt.config.Period)

	loop:
		for {
			err, sfr := bt.GetSafecast(b, lastTime)
			if err != nil {
				bt.logger.Errorf("Error while getting Safecast data: %v", err)
			} else {
				if sfr.success {
					lastTime = sfr.tm
				}
			}

			select {
			case <-bt.done:
				break loop
			case <-ticker.C:
			}
		}
	}()

	<-bt.done

	return nil
}

// Stop stops safecastbeat.
func (bt *Safecastbeat) Stop() {
	if err := bt.client.Close(); err != nil {
		bt.logger.Warnf("Error closing client: %v", err)
	}
	close(bt.done)
}

type SFResponse struct {
	tm      time.Time
	success bool
}

func (bt *Safecastbeat) GetSafecast(b *beat.Beat, lastRun time.Time) (error, SFResponse) {
	now := time.Now().UTC()
	sfr := SFResponse{tm: now, success: false}

	parsedURL, err := url.Parse(safecastUrl)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err), sfr
	}

	parsedURL.Path += safecastPath
	parameters := url.Values{}
	parameters.Add("per_page", perPage)
	parameters.Add("since", lastRun.Format(dateTimeFormat))
	parsedURL.RawQuery = parameters.Encode()

	bt.logger.Debugf("Last run at: %s", lastRun.String())
	bt.logger.Debugf("Requesting Safecast data: %s", parsedURL.String())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-bt.done:
			cancel()
		case <-ctx.Done():
		}
	}()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err), sfr
	}

	res, err := bt.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err), sfr
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bt.logger.Debugf("Unexpected status code: %d", res.StatusCode)
		return fmt.Errorf("HTTP %d from %s", res.StatusCode, parsedURL.String()), sfr
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		bt.logger.Errorf("Failed to read response body: %v", err)
		return fmt.Errorf("reading response body: %w", err), sfr
	}
	bt.logger.Debugf("Response body size: %d bytes", len(body))

	sfdata := []SafecastResponseData{}
	err = json.Unmarshal(body, &sfdata)
	if err != nil {
		bt.logger.Errorf("Failed to unmarshal JSON: %v", err)
		return fmt.Errorf("unmarshalling JSON: %w", err), sfr
	}

	if len(sfdata) > 0 {
		bt.logger.Debugf("%d new documents.", len(sfdata))

		transformedData := bt.TransformStationData(sfdata)

		ts := time.Now()
		for _, d := range transformedData {
			event := beat.Event{
				Timestamp: ts,
				Fields: mapstr.M{
					"type":     "safecastbeat",
					"safecast": d,
				},
			}
			bt.client.Publish(event)
		}
		sfr.success = true
	} else {
		sfr.success = true
		bt.logger.Debug("No new data.")
	}

	return nil, sfr
}

func (bt *Safecastbeat) TransformStationData(data []SafecastResponseData) []mapstr.M {
	safecastData := []mapstr.M{}

	for _, d := range data {
		sf := mapstr.M{
			"safecastId":          d.Id,
			"userId":              d.User_id,
			"value":               d.Value,
			"unit":                d.Unit,
			"locationName":        d.Location_name,
			"deviceId":            d.Device_id,
			"originalId":          d.Original_id,
			"measurementImportId": d.Measurement_import_id,
			"capturedAt":          d.Captured_at,
			"height":              d.Height,
			"devicetypeId":        d.Devicetype_id,
			"sensorId":            d.Sensor_id,
			"stationId":           d.Station_id,
			"channelId":           d.Channel_id,
			"location": mapstr.M{
				"lat": d.Latitude,
				"lon": d.Longitude,
			},
		}

		safecastData = append(safecastData, sf)
	}

	return safecastData
}
