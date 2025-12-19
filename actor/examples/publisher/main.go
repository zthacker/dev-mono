package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// MeasurementPacket matches the type in estimation/types.go
type MeasurementPacket struct {
	SatelliteID int32   `json:"satellite_id"`
	Timestamp   int64   `json:"timestamp"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	PosZ        float64 `json:"pos_z"`
	CovXX       float64 `json:"cov_xx"`
	CovYY       float64 `json:"cov_yy"`
	CovZZ       float64 `json:"cov_zz"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting test publisher...")

	// Connect to NATS
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure stream exists
	ctx := context.Background()
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "TELEMETRY",
		Description: "Satellite telemetry measurements",
		Subjects:    []string{"telemetry.>"},
		Retention:   jetstream.WorkQueuePolicy,
		MaxAge:      24 * time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Publishing test measurements...")

	// Simulate measurements for 10 satellites across 2 partitions
	numSatellites := 10
	numPartitions := 2
	messagesPerSat := 5

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < messagesPerSat; i++ {
		for satID := 0; satID < numSatellites; satID++ {
			// Generate simulated GPS measurement
			// (In reality, this would be orbital position in ECI)
			measurement := MeasurementPacket{
				SatelliteID: int32(satID),
				Timestamp:   time.Now().UnixNano(),
				// Simulate LEO orbit (~7000 km altitude in ECI)
				PosX:  7000000.0 + rand.Float64()*1000.0,
				PosY:  rand.Float64() * 1000000.0,
				PosZ:  rand.Float64() * 1000000.0,
				CovXX: 100.0, // 100 m^2 measurement uncertainty
				CovYY: 100.0,
				CovZZ: 100.0,
			}

			data, err := json.Marshal(measurement)
			if err != nil {
				log.Printf("Failed to marshal: %v", err)
				continue
			}

			// Route to correct partition (sharding)
			partition := (satID % numPartitions) + 1
			subject := fmt.Sprintf("telemetry.identified.partition.%d.%d", partition, satID)

			if _, err := js.Publish(ctx, subject, data); err != nil {
				log.Printf("Failed to publish: %v", err)
				continue
			}

			log.Printf("Published: sat=%d partition=%d pos=[%.0f, %.0f, %.0f]",
				satID, partition, measurement.PosX, measurement.PosY, measurement.PosZ)
		}

		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Published %d measurements across %d satellites", messagesPerSat*numSatellites, numSatellites)
	log.Println("Done!")
}
