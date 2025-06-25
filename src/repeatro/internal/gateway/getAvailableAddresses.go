package gateway

import (
	"context"
	"math/rand"
)

func GetAvailableAddresses(ctx context.Context, serviceID, endpoint string, serviceAddresses func(ctx context.Context, serviceID string) ([]string, error)) (string, error) {
	addrs, err := serviceAddresses(ctx, serviceID)
	if err != nil {
		return "", err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + endpoint
	return url, nil
}