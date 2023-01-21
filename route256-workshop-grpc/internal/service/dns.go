package service

import (
	"context"
	"sync"
)

type DNS struct {
	addressBindings map[string][]string
	abmx            sync.RWMutex
}

func NewInMemDNS() *DNS {
	return &DNS{
		addressBindings: map[string][]string{},
	}
}

func (d *DNS) Register(
	ctx context.Context,
	serviceName string,
	addresses ...string,
) {
	d.abmx.Lock()
	defer d.abmx.Unlock()

	d.addressBindings[serviceName] = uniqAddresses(
		append(d.addressBindings[serviceName], addresses...)...,
	)
}

func (d *DNS) Unregister(
	ctx context.Context,
	serviceName string,
	addresses ...string,
) {
	d.abmx.Lock()
	defer d.abmx.Unlock()

	pool, ok := d.addressBindings[serviceName]
	if !ok {
		return
	}

	for _, address := range addresses {
		pool = d.removeAddressForService(pool, address)
	}

	d.addressBindings[serviceName] = pool
}

func (d *DNS) removeAddressForService(
	pool []string,
	address string,
) []string {
	for i := range pool {
		if address == pool[i] {
			pool[i] = pool[len(pool)-1]
			pool[len(pool)-1] = ""
			pool = pool[:len(pool)-1]
		}
	}

	return pool
}

func (d *DNS) GetAddressesForService(
	ctx context.Context,
	serviceName string,
) ([]string, bool) {
	d.abmx.RLock()
	defer d.abmx.RUnlock()

	addresses, arePresent := d.addressBindings[serviceName]
	return addresses, arePresent
}

func uniqAddresses(addresses ...string) []string {
	addressesMap := make(map[string]any, len(addresses))
	for _, v := range addresses {
		// Активно пользуемся тем, что пустые структуры удовлетворяют any и при этом ничего не весят
		addressesMap[v] = struct{}{}
	}

	uniqueAddresses := []string{}
	for k := range addressesMap {
		uniqueAddresses = append(uniqueAddresses, k)
	}

	return uniqueAddresses
}
