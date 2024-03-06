package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	TestEvent = &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mims-mas-product-service-5ff5668cdf-qppr2.17ba228424da4149",
		},
		Type: "Warning",
		InvolvedObject: v1.ObjectReference{
			Kind:      "Node",
			Namespace: "default",
		},
		Reason: "BackOff",
	}
)

func TestEvents(t *testing.T) {
	kindFilter := NewGenericFilter("Kind", []string{"Node"}, false)
	assert.True(t, kindFilter.Filter(TestEvent), "")

	typeFilter := NewGenericFilter("Type", []string{"Warning"}, false)
	assert.True(t, typeFilter.Filter(TestEvent), "")

	namespaceFilter := NewGenericFilter("Namespace", []string{"default"}, false)
	assert.True(t, namespaceFilter.Filter(TestEvent), "")

	reasonFilter := NewGenericFilter("Reason", []string{"BackOff"}, false)
	assert.True(t, reasonFilter.Filter(TestEvent), "")

	regexReasonFilter := NewGenericFilter("Reason", []string{"BackOff"}, true)
	assert.True(t, regexReasonFilter.Filter(TestEvent), "")

	regexReasonsFilter := NewGenericFilter("Reason", []string{"Unhealthy", "BackOff"}, true)
	assert.True(t, regexReasonsFilter.Filter(TestEvent), "")

	objectFilter := NewGenericFilter("Object", []string{"mims-mas-product-service"}, true)
	assert.True(t, objectFilter.Filter(TestEvent))
}
