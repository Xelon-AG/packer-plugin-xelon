package builder

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/useragent"

	"github.com/Xelon-AG/packer-plugin-xelon/internal/version"
	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

func newXelonClient(c Config) (*xelon.Client, error) {
	opts := []xelon.ClientOption{xelon.WithUserAgent(useragent.String(version.PluginVersion.FormattedVersion()))}
	if c.BaseURL != "" {
		opts = append(opts, xelon.WithBaseURL(c.BaseURL))
	}
	opts = append(opts, xelon.WithClientID(c.ClientID))

	return xelon.NewClient(c.Token, opts...), nil
}

const (
	maxRetries   = 10
	baseDelay    = 5 * time.Second
	maxDelay     = 30 * time.Second
	pollInterval = 15 * time.Second
)

func waitDevicePowerStateOn(ctx context.Context, client *xelon.Client, deviceID string) error {
	attempt := 0

	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	log.Printf("[INFO] Start waiting for device (%s) to be powered on", deviceID)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for device to be powered on cancelled: %w", ctx.Err())
		default:
			attempt++

			device, _, err := client.Devices.Get(ctx, deviceID)
			if err != nil {
				if attempt >= maxRetries {
					return fmt.Errorf("failed to wait for device to be powered on after %d attempts: %w", attempt, err)
				}
				delay := exponentialBackoff(attempt, baseDelay, maxDelay)
				log.Printf("[DEBUG] API error (attempt %d/%d). Retrying in %d: %v", attempt, maxRetries, delay, err)
				time.Sleep(delay)
			}

			log.Printf("[INFO] Device info - id: %s, poweredOn: %v, state: %d", device.ID, device.PoweredOn, device.State)

			if device.PoweredOn {
				log.Printf("[INFO] Device (%s) is powered on", deviceID)
				return nil
			}

			if attempt >= maxRetries {
				return fmt.Errorf("timeout: device (%s) did not reach powered on state after %d attempts", deviceID, attempt)
			}

			time.Sleep(pollInterval)
		}
	}
}

func waitDevicePowerStateOff(ctx context.Context, client *xelon.Client, deviceID string) error {
	attempt := 0

	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	log.Printf("[INFO] Start waiting for device (%s) to be powered off", deviceID)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for device to be powered off cancelled: %w", ctx.Err())
		default:
			attempt++

			device, _, err := client.Devices.Get(ctx, deviceID)
			if err != nil {
				if attempt >= maxRetries {
					return fmt.Errorf("failed to wait for device to be powered off after %d attempts: %w", attempt, err)
				}
				delay := exponentialBackoff(attempt, baseDelay, maxDelay)
				log.Printf("[DEBUG] API error (attempt %d/%d). Retrying in %d: %v", attempt, maxRetries, delay, err)
				time.Sleep(delay)
			}

			log.Printf("[DEBUG] Device info - id: %s, poweredOn: %v, state: %d", device.ID, device.PoweredOn, device.State)

			if !device.PoweredOn {
				log.Printf("[INFO] Device (%s) is powered off", deviceID)
				return nil
			}

			if attempt >= maxRetries {
				return fmt.Errorf("timeout: device (%s) did not reach powered off state after %d attempts", deviceID, attempt)
			}

			time.Sleep(pollInterval)
		}
	}
}

func waitDeviceStateReady(ctx context.Context, client *xelon.Client, deviceID string) error {
	attempt := 0

	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	log.Printf("[INFO] Start waiting for device (%s) to be ready", deviceID)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for device to be ready cancelled: %w", ctx.Err())
		default:
			attempt++

			device, _, err := client.Devices.Get(ctx, deviceID)
			if err != nil {
				if attempt >= maxRetries {
					return fmt.Errorf("failed to wait for device to be ready after %d attempts: %w", attempt, err)
				}
				delay := exponentialBackoff(attempt, baseDelay, maxDelay)
				log.Printf("[DEBUG] API error (attempt %d/%d). Retrying in %d: %v", attempt, maxRetries, delay, err)
				time.Sleep(delay)
			}

			log.Printf("[DEBUG] Device info - id: %s, poweredOn: %v, state: %d", device.ID, device.PoweredOn, device.State)

			if device.State == 1 {
				log.Printf("[INFO] Device (%s) is ready", deviceID)
				return nil
			}

			if attempt >= maxRetries {
				return fmt.Errorf("timeout: device (%s) did not reach ready state after %d attempts", deviceID, attempt)
			}

			time.Sleep(pollInterval)
		}
	}
}

func exponentialBackoff(attempt int, base, max time.Duration) time.Duration {
	delay := base * time.Duration(1<<(attempt-1))
	if delay > max {
		return max
	}
	jitter := time.Duration(rand.Int63n(int64(base)))
	return delay + jitter
}
