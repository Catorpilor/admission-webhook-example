package main

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/klog"
)

// This function expects all CRDs submitted to it to be apiextensions.k8s.io/v1beta1
// TODO: When apiextensions.k8s.io/v1 is added we will need to update this function.
func admitCRD(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("admitting crd")
	crdResource := metav1.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	if ar.Request.Resource != crdResource {
		err := fmt.Errorf("expect resource to be %s", crdResource)
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	crd := apiextensionsv1.CustomResourceDefinition{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &crd); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true

	if v, ok := crd.Labels["webhook-e2e-test"]; ok {
		if v == "webhook-disallow" {
			reviewResponse.Allowed = false
			reviewResponse.Result = &metav1.Status{Message: "the crd contains unwanted label"}
		}
	}
	return &reviewResponse
}
