package main

import (
	"encoding/json"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	customResourcePatch1 string = `[
         { "op": "add", "path": "/data/mutation-stage-1", "value": "yes" }
     ]`
	customResourcePatch2 string = `[
         { "op": "add", "path": "/data/mutation-stage-2", "value": "yes" }
     ]`
)

func mutateCustomResource(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("mutating custom resource")
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true

	if cr.Data["mutation-start"] == "yes" {
		reviewResponse.Patch = []byte(customResourcePatch1)
	}
	if cr.Data["mutation-stage-1"] == "yes" {
		reviewResponse.Patch = []byte(customResourcePatch2)
	}
	pt := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse
}

func admitCustomResource(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("admitting custom resource")
	cr := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true
	for k, v := range cr.Data {
		if k == "webhook-e2e-test" && v == "webhook-disallow" {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the custom resource contains unwanted data",
			}
		}
	}
	return &reviewResponse
}
