package service

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	calcv1alpha1 "github.com/mykysha/kubCalculator/api/v1alpha1"
)

type Repository interface {
	// ProcessCalculator processes the given calculator.
	ProcessCalculator(ctx context.Context, calc *calcv1alpha1.Calculator) error
	// DefineSecret defines a calculator operator secret.
	DefineSecret(ctx context.Context, name, namespace string, result int) (*corev1.Secret, error)
}

type CalculatorService struct{}

// ProcessCalculator processes the given calculator.
func (c CalculatorService) ProcessCalculator(_ context.Context, calc *calcv1alpha1.Calculator) error {
	calc.Status.Result = calc.Spec.X + calc.Spec.Y
	calc.Status.Processed = true

	return nil
}

// DefineSecret defines a calculator operator secret.
func (c CalculatorService) DefineSecret(_ context.Context, name, namespace string, result int) (*corev1.Secret, error) {
	data := make(map[string]string)
	data["result"] = strconv.Itoa(result)

	annotations := make(map[string]string)
	annotations["managed-by"] = "calc-operator"

	immutable := false

	return &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Annotations: annotations},
		Immutable:  &immutable,
		Data:       map[string][]byte{},
		StringData: data,
		Type:       "Opaque",
	}, nil
}
