package controllers

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	calcv1alpha1 "github.com/mykysha/kubCalculator/api/v1alpha1"
	mocks "github.com/mykysha/kubCalculator/mocks/pkg/service"
)

var _ = Describe("Calculator controller", func() {
	Context("Calculator controller test", func() {
		const CalculatorName = "test-reconciler"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      CalculatorName,
				Namespace: CalculatorName,
			},
		}

		typeNamespaceName := types.NamespacedName{
			Name:      CalculatorName,
			Namespace: CalculatorName,
		}

		BeforeEach(func() {
			By("Creating namespace")

			err := k8sClient.Create(ctx, namespace)
			Expect(err).NotTo(HaveOccurred())

			By("Set image env")
			err = os.Setenv("OPERATOR_IMAGE", "nndergunov/kubcalculator:latest")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			By("Deleting namespace")

			err := k8sClient.Delete(ctx, namespace)
			Expect(err).NotTo(HaveOccurred())

			By("Unset image env")
			err = os.Unsetenv("OPERATOR_IMAGE")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should successfully reconcile", func() {
			By("Creating the custom resource for the Kind Memcached")
			calculator := &calcv1alpha1.Calculator{}
			err := k8sClient.Get(ctx, typeNamespaceName, calculator)
			if err != nil && errors.IsNotFound(err) {
				calculator := &calcv1alpha1.Calculator{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Calculator",
						APIVersion: "calc.example.com/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      CalculatorName,
						Namespace: CalculatorName,
					},
					Spec: calcv1alpha1.CalculatorSpec{
						X: 2,
						Y: 3,
					},
				}

				err = k8sClient.Create(ctx, calculator)
				Expect(err).To(Not(HaveOccurred()))
			}

			By("Checking if the custom resource was successfully created")
			found := &calcv1alpha1.Calculator{}
			err = k8sClient.Get(ctx, typeNamespaceName, found)
			Expect(err).To(Not(HaveOccurred()))

			found.Status.Result = found.Spec.X + found.Spec.Y
			found.Status.Processed = true

			var repo mocks.Repository

			repo.On("ProcessCalculator", ctx, mock.AnythingOfType("*v1alpha1.Calculator")).Return(
				nil).Run(func(args mock.Arguments) {
				returnCalc := new(calcv1alpha1.Calculator)

				switch c := args.Get(1).(type) {
				case *calcv1alpha1.Calculator:
					returnCalc = c
				default:
					Fail(fmt.Sprintf("Unexpected type: %T", c))
				}

				returnCalc.Status = found.Status
			})

			repo.On("DefineSecret", ctx, mock.AnythingOfType("string"),
				mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(&corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        CalculatorName,
					Namespace:   CalculatorName,
					Annotations: map[string]string{"managed-by": "calc-operator"},
				},
				Immutable:  new(bool),
				StringData: map[string]string{"result": fmt.Sprintf("%d", 5)},
				Type:       "Opaque",
			}, nil)

			memcachedReconciler := &CalculatorReconciler{
				Client:  k8sClient,
				Service: &repo,
				Scheme:  k8sClient.Scheme(),
			}

			By("Create manager for reconciler")
			mgr, err := ctrl.NewManager(cfg, ctrl.Options{
				Scheme: k8sClient.Scheme(),
			})
			Expect(err).ToNot(HaveOccurred())

			By("Setup reconciler with manager")
			err = memcachedReconciler.SetupWithManager(mgr, &repo)
			Expect(err).ToNot(HaveOccurred())

			_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if the latest Status processed is true, result is x+y and is added to the Calculator instance")
			calculator = &calcv1alpha1.Calculator{}
			err = k8sClient.Get(ctx, typeNamespaceName, calculator)
			Expect(err).To(Not(HaveOccurred()))
			if calculator.Status.Processed != true || calculator.Status.Result != calculator.Spec.X+calculator.Spec.Y {
				Fail("Status is not updated")
			}
		})
	})
})

var _ = Describe("Secret manager", func() {
	Context("Secret manager test", func() {
		const CalculatorName = "test-secret-manager"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      CalculatorName,
				Namespace: CalculatorName,
			},
		}

		typeNamespaceName := types.NamespacedName{
			Name:      CalculatorName,
			Namespace: CalculatorName,
		}

		BeforeEach(func() {
			By("Creating namespace")

			err := k8sClient.Create(ctx, namespace)
			Expect(err).NotTo(HaveOccurred())

			By("Set image env")
			err = os.Setenv("OPERATOR_IMAGE", "nndergunov/kubcalculator:latest")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			By("Deleting namespace")

			err := k8sClient.Delete(ctx, namespace)
			Expect(err).NotTo(HaveOccurred())

			By("Unset image env")
			err = os.Unsetenv("OPERATOR_IMAGE")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should update secret if it exists", func() {
			By("Creating the secret if it does not exist")
			secret := &corev1.Secret{}
			err := k8sClient.Get(ctx, typeNamespaceName, secret)
			if err != nil && errors.IsNotFound(err) {
				secret := &corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Secret",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:        CalculatorName,
						Namespace:   CalculatorName,
						Annotations: map[string]string{"managed-by": "calc-operator"},
					},
					Immutable:  new(bool),
					StringData: map[string]string{"result": fmt.Sprintf("%d", 5)},
					Type:       "Opaque",
				}

				err = k8sClient.Create(ctx, secret)
				Expect(err).To(Not(HaveOccurred()))
			}

			By("Checking if the secret was successfully created")
			found := &corev1.Secret{}
			err = k8sClient.Get(ctx, typeNamespaceName, found)
			Expect(err).To(Not(HaveOccurred()))

			if found.Data != nil {
				found.Data["result"] = []byte(fmt.Sprintf("%d", 10))
			} else if found.StringData != nil {
				found.StringData["result"] = fmt.Sprintf("%d", 10)
			} else {
				Fail("Secret data is nil")
			}

			reconciler := &CalculatorReconciler{
				Client:  k8sClient,
				Service: nil,
				Scheme:  k8sClient.Scheme(),
			}

			By("Manage updated secret")
			err = reconciler.manageSecret(ctx, found)
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if the secret was successfully updated")
			found = &corev1.Secret{}
			err = k8sClient.Get(ctx, typeNamespaceName, found)
			Expect(err).To(Not(HaveOccurred()))
			if found.Data != nil {
				if string(found.Data["result"]) != "10" {
					Fail("Secret data is not updated")
				}
			} else if found.StringData != nil {
				if found.StringData["result"] != "10" {
					Fail("Secret data is not updated")
				}
			} else {
				Fail("Updated secret data is nil")
			}
		})
	})
})
