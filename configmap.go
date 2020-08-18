package main

import (
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	configMapPatch1 string = `[
         { "op": "add", "path": "/data/mutation-stage-1", "value": "yes" }
     ]`
	configMapPatch2 string = `[
         { "op": "add", "path": "/data/mutation-stage-2", "value": "yes" }
     ]`
)

// deny configmaps with specific key-value pair.
func admitConfigMaps(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("admitting configmaps")
	configMapResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	if ar.Request.Resource != configMapResource {
		klog.Errorf("expect resource to be %s", configMapResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	configmap := corev1.ConfigMap{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &configmap); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true
	for k, v := range configmap.Data {
		if k == "webhook-e2e-test" && v == "webhook-disallow" {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{
				Reason: "the configmap contains unwanted key and value",
			}
		}
	}
	return &reviewResponse
}

func mutateConfigmaps(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("mutating configmaps")
	configMapResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	if ar.Request.Resource != configMapResource {
		klog.Errorf("expect resource to be %s", configMapResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	configmap := corev1.ConfigMap{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &configmap); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if configmap.Data["mutation-start"] == "yes" {
		reviewResponse.Patch = []byte(configMapPatch1)
	}
	if configmap.Data["mutation-stage-1"] == "yes" {
		reviewResponse.Patch = []byte(configMapPatch2)
	}

	pt := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return &reviewResponse
}
