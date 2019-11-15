package beater

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/radoondas/safecastbeat/config"
)

// Safecastbeat configuration.
type Safecastbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
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

var lastRunAdded, lastRunCaptured, lastRunPoint time.Time

// New creates an instance of safecastbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Safecastbeat{
		done:   make(chan struct{}),
		config: c,
	}

	err := bt.init(b)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

// Run starts safecastbeat.
func (bt *Safecastbeat) Run(b *beat.Beat) error {
	logp.NewLogger(selector).Info("safecastbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	logp.NewLogger(selector).Info("Configuration: ", bt.config)

	go func() {
		ticker := time.NewTicker(bt.config.Period)
		defer ticker.Stop()

		// UTC is required
		lastTime := time.Now().UTC().Add(-bt.config.Period)

		for {

			err, sfr := bt.GetSafecast(b, lastTime, bt.config.SSL.CAfile)
			if err != nil {
				logp.NewLogger(selector).Error("Error while getting Safecast data: %v", err)
			} else {
				if sfr.success {
					lastTime = sfr.tm
				}
			}

			select {
			case <-bt.done:
				goto GotoFinish
			case <-ticker.C:
			}
		}
	GotoFinish:
	}()

	<-bt.done

	return nil
}

// Stop stops safecastbeat.
func (bt *Safecastbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Safecastbeat) init(b *beat.Beat) error {

	// Configuration checks
	if bt.config.SSL.CAfile == "" {
		logp.NewLogger(selector).Info("CAFile IS NOT set.")
	}

	return nil
}

type SFResponse struct {
	tm      time.Time
	success bool
}

func (bt *Safecastbeat) GetSafecast(b *beat.Beat, lastRun time.Time, CAFile string) (error, SFResponse) {

	now := time.Now().UTC()
	sfr := SFResponse{tm: now, success: false}

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	var ParsedUrl *url.URL

	if CAFile != "" {
		// Load our trusted certificate path
		pemData, err := ioutil.ReadFile(CAFile)
		if err != nil {
			panic(err)
		}
		ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
		if !ok {
			logp.NewLogger(selector).Error("Unable to load CA file")
			panic("Couldn't load PEM data")
		}
	}

	client := &http.Client{Transport: transport}

	ParsedUrl, err := url.Parse(safecastUrl)
	if err != nil {
		logp.NewLogger(selector).Error("Unable to parse URL string")
		panic(err)
	}

	ParsedUrl.Path += safecastPath
	parameters := url.Values{}

	parameters.Add("per_page", perPage)
	parameters.Add("since", lastRun.Format(dateTimeFormat))
	logp.NewLogger(selector).Debug("Last run at: ", lastRun.String())

	ParsedUrl.RawQuery = parameters.Encode()
	logp.NewLogger(selector).Debug("Requesting Safecast data: ", ParsedUrl.String())

	req, err := http.NewRequest("GET", ParsedUrl.String(), nil)

	res, err := client.Do(req)

	if err != nil {
		return err, sfr
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logp.NewLogger(selector).Debug("Status code: ", res.StatusCode)
		logp.NewLogger(selector).Debug("Status code: ", res.Body)
		return fmt.Errorf("HTTP %+v", res), sfr
	}

	body, err := ioutil.ReadAll(res.Body)
	logp.NewLogger(selector).Debug(body)
	if err != nil {
		log.Fatal(err)
		return err, sfr
	}

	sfdata := []SafecastResponseData{}
	err = json.Unmarshal([]byte(body), &sfdata)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err, sfr
	}

	if len(sfdata) > 0 {
		logp.NewLogger(selector).Debug(len(sfdata), " new documents.")
		//logp.NewLogger(selector).Debug("Unmarshal-ed Safecast data: ", sfdata)

		transformedData := bt.TransformStationData(sfdata)

		ts := time.Now()
		for _, d := range transformedData {

			event := beat.Event{
				Timestamp: ts,
				Fields: common.MapStr{
					"type":     "safecastbeat",
					"safecast": d,
				},
			}

			bt.client.Publish(event)
			logp.NewLogger(selector).Debug("Event: ", event)
		}
		sfr.success = true

	} else {
		sfr.success = true
		logp.NewLogger(selector).Debug("No new data.")
	}

	return nil, sfr
}

func (bt *Safecastbeat) TransformStationData(data []SafecastResponseData) []common.MapStr {

	safecastData := []common.MapStr{}

	for _, d := range data {
		//logp.NewLogger(selector).Debug("Device data: ", device)

		sf := common.MapStr{
			"safecastId":          d.Id,
			"userId":              d.User_id,
			"value":               d.Value,
			"unit":                d.Unit,
			"locationName":        d.Location_name,
			"deviceId":            d.Device_id,
			"originalId":          d.Original_id,
			"measurementImportId": d.Measurement_import_id,
			"capturedAt":          d.Captured_at, //date
			"height":              d.Height,
			"devicetypeId":        d.Devicetype_id,
			"sensorId":            d.Sensor_id,
			"stationId":           d.Station_id,
			"channelId":           d.Channel_id,
			"location": common.MapStr{
				"lat": d.Latitude,
				"lon": d.Longitude,
			},
		}
		logp.NewLogger(selector).Debug("SF data: ", sf)

		safecastData = append(safecastData, sf)
	}

	return safecastData
}
