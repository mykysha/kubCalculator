package service_test

import (
	"context"
	"testing"

	calcv1alpha1 "github.com/mykysha/kubCalculator/api/v1alpha1"
	"github.com/mykysha/kubCalculator/pkg/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProcessCalculator(t *testing.T) { //nolint:funlen // Test function
	t.Parallel()

	// Test table
	tests := []struct {
		name     string
		calc     *calcv1alpha1.Calculator
		expected *calcv1alpha1.Calculator
	}{
		{
			name: "Max int32 x & y no status",
			calc: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 2147483647,
					Y: 2147483647,
				},
			},
			expected: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 2147483647,
					Y: 2147483647,
				},
				Status: calcv1alpha1.CalculatorStatus{
					Processed: true,
					Result:    4294967294,
				},
			},
		},
		{
			name: "Reprocess calculator",
			calc: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 1,
					Y: 1,
				},
				Status: calcv1alpha1.CalculatorStatus{
					Processed: true,
					Result:    0,
				},
			},
			expected: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 1,
					Y: 1,
				},
				Status: calcv1alpha1.CalculatorStatus{
					Processed: true,
					Result:    2,
				},
			},
		},
		{
			name: "Correct result, processed status is false",
			calc: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 1,
					Y: 1,
				},
				Status: calcv1alpha1.CalculatorStatus{
					Processed: false,
					Result:    2,
				},
			},
			expected: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 1,
					Y: 1,
				},
				Status: calcv1alpha1.CalculatorStatus{
					Processed: true,
					Result:    2,
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := &service.CalculatorService{}
			ctx := context.Background()

			// Process calculator
			err := s.ProcessCalculator(ctx, test.calc)
			if err != nil {
				t.Fatalf("processCalculator() error = %v", err)
			}

			// Check result
			if test.calc.Status.Result != test.expected.Status.Result {
				t.Errorf("processCalculator() got = %v, want %v", test.calc.Status.Result, test.expected.Status.Result)
			}

			if test.calc.Status.Processed != test.expected.Status.Processed {
				t.Errorf("processCalculator() got = %v, want %v", test.calc.Status.Processed, test.expected.Status.Processed)
			}

			if test.calc.Spec.X != test.expected.Spec.X {
				t.Errorf("processCalculator() got = %v, want %v", test.calc.Spec.X, test.expected.Spec.X)
			}

			if test.calc.Spec.Y != test.expected.Spec.Y {
				t.Errorf("processCalculator() got = %v, want %v", test.calc.Spec.Y, test.expected.Spec.Y)
			}
		})
	}
}

func TestDefineSecret(t *testing.T) { //nolint:funlen // Test function
	t.Parallel()

	// Test table
	tests := []struct {
		name       string
		result     int
		secretName string
		namespace  string
		expected   *corev1.Secret
	}{
		{
			name:       "Define secret",
			result:     2,
			secretName: "test-secret",
			namespace:  "non-default",
			expected: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-secret",
					Namespace:   "non-default",
					Annotations: map[string]string{"managed-by": "calc-operator"},
				},
				StringData: map[string]string{
					"result": "2",
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := &service.CalculatorService{}
			ctx := context.Background()

			// Define secret
			got, err := s.DefineSecret(ctx, test.secretName, test.namespace, test.result)
			if err != nil {
				t.Fatalf("defining secret error = %v", err)
			}

			// Check result
			if got.Name != test.expected.Name {
				t.Errorf("secret name: got = %v, want = %v", got.Name, test.expected.Name)
			}

			if got.Namespace != test.expected.Namespace {
				t.Errorf("secret namespace: got = %v, want = %v", got.Namespace, test.expected.Namespace)
			}

			if got.StringData["result"] != test.expected.StringData["result"] {
				t.Errorf("secret result: got = %v, want = %v", got.Data["result"], test.expected.Data["result"])
			}

			if got.Annotations["managed-by"] != test.expected.Annotations["managed-by"] {
				t.Errorf("secret annotations: got = %v, want = %v",
					got.Annotations["managed-by"], test.expected.Annotations["managed-by"])
			}
		})
	}
}
