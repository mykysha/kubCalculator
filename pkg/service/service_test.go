package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	calcv1alpha1 "github.com/mykysha/kubCalculator/api/v1alpha1"
	"github.com/mykysha/kubCalculator/pkg/service"
)

func TestProcessCalculator(t *testing.T) { //nolint:funlen // Test function
	t.Parallel()

	// Test table
	tests := []struct {
		name string
		calc *calcv1alpha1.Calculator
		want *calcv1alpha1.Calculator
	}{
		{
			name: "Max int32 x & y no status",
			calc: &calcv1alpha1.Calculator{
				Spec: calcv1alpha1.CalculatorSpec{
					X: 2147483647,
					Y: 2147483647,
				},
			},
			want: &calcv1alpha1.Calculator{
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
			want: &calcv1alpha1.Calculator{
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
			want: &calcv1alpha1.Calculator{
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &service.CalculatorService{}
			ctx := context.Background()

			// Process calculator
			assert.NoError(t, s.ProcessCalculator(ctx, tt.calc))

			// Check result
			assert.Equal(t, tt.want, tt.calc)
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
		want       *corev1.Secret
	}{
		{
			name:       "Define secret",
			result:     2,
			secretName: "test-secret",
			namespace:  "non-default",
			want: &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-secret",
					Namespace:   "non-default",
					Annotations: map[string]string{"managed-by": "calc-operator"},
				},
				Immutable: pointer.Bool(false),
				Data:      make(map[string][]uint8),
				StringData: map[string]string{
					"result": "2",
				},
				Type: "Opaque",
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &service.CalculatorService{}
			ctx := context.Background()

			// Define secret
			got, err := s.DefineSecret(ctx, tt.secretName, tt.namespace, tt.result)
			assert.NoError(t, err)

			// Check result
			assert.Equal(t, tt.want, got)
		})
	}
}
