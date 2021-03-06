// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/event"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
)

// GetJobEvents gets events associated to job.
func GetJobEvents(client client.Interface, namespace, jobName string) (
	*common.EventList, error) {

	log.Printf("Getting events related to %s job in %s namespace", jobName,
		namespace)

	// Get events for job.
	rsEvents, err := event.GetEvents(client, namespace, jobName)

	if err != nil {
		return nil, err
	}

	// Get events for pods in job.
	podEvents, err := GetJobPodsEvents(client, namespace, jobName)

	if err != nil {
		return nil, err
	}

	apiEvents := append(rsEvents, podEvents...)

	if !event.IsTypeFilled(apiEvents) {
		apiEvents = event.FillEventsType(apiEvents)
	}

	// TODO support pagination
	events := event.CreateEventList(apiEvents, dataselect.NoDataSelect)

	log.Printf("Found %d events related to %s job in %s namespace",
		len(events.Events), jobName, namespace)

	return &events, nil
}

// GetJobPodsEvents gets events associated to pods in job.
func GetJobPodsEvents(client client.Interface, namespace, jobName string) (
	[]api.Event, error) {

	job, err := client.Extensions().Jobs(namespace).Get(jobName)

	if err != nil {
		return nil, err
	}

	podEvents, err := event.GetPodsEvents(client, namespace, job.Spec.Selector.MatchLabels)

	if err != nil {
		return nil, err
	}

	return podEvents, nil
}
