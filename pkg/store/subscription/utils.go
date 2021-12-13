// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package subscription

import (
	"context"
	"fmt"
)

func SubscriptionsByType(ctx context.Context, s Store, subscriptionType int, ch chan<- *Entry) error {

	subscriptionTypeString := SubscriptionType(subscriptionType).String()
	subchEntries := make(chan *Entry)
	done := make(chan bool)

	go func(subscriptionTypeString string, ch chan<- *Entry) {

		for entry := range subchEntries {
			subValue := entry.Value.(Value)
			for _, subEntry := range subValue.Subscriptions {
				if subEntry.Type == subscriptionTypeString {
					ch <- entry
				}
			}
		}
		done <- true
	}(subscriptionTypeString, ch)

	err := s.Entries(ctx, subchEntries)
	if err != nil {
		close(ch)
		return fmt.Errorf("no subscriptions entries stored for Type %s", subscriptionTypeString)
	}

	<-done
	close(ch)
	return nil
}

func SubscriptionsByTypeID(ctx context.Context, s Store, subscriptionType int, subscriptionTypeID string, ch chan<- *Entry) error {

	subscriptionTypeString := SubscriptionType(subscriptionType).String()
	subchEntries := make(chan *Entry)
	done := make(chan bool)

	go func(subscriptionTypeString, subscriptionTypeID string, ch chan<- *Entry) {

		for entry := range subchEntries {
			subValue := entry.Value.(Value)
			for _, subEntry := range subValue.Subscriptions {
				if subEntry.Type == subscriptionTypeString && subEntry.TypeID == subscriptionTypeID {
					ch <- entry
				}
			}
		}
		done <- true
	}(subscriptionTypeString, subscriptionTypeID, ch)

	err := s.Entries(ctx, subchEntries)
	if err != nil {
		close(ch)
		return fmt.Errorf("no subscription entries stored for TypeID %s", subscriptionTypeID)
	}

	<-done
	close(ch)
	return nil
}
