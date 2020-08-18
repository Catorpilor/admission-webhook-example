package main

import (
	"encoding/json"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	addFirstLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels", "value": {"added-label": "yes"}}
     ]`
	addAdditionalLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels/added-label", "value": "yes" }
     ]`
)

// Add a label {"added-label": "yes"} to the object
func addLabel(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("calling add-label")
	obj := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}
	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if len(obj.ObjectMeta.Labels) == 0 {
		reviewResponse.Patch = []byte(addFirstLabelPatch)
	} else {
		reviewResponse.Patch = []byte(addAdditionalLabelPatch)
	}
	pt := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse
}
