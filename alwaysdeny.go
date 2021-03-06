package main

import (
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// alwaysDeny all requests made to this function.
func alwaysDeny(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("calling always-deny")
	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = false
	reviewResponse.Result = &metav1.Status{Message: "this webhook denies all requests"}
	return &reviewResponse
}
