/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	calcv1alpha1 "github.com/mykysha/kubCalculator/api/v1alpha1"
	"github.com/mykysha/kubCalculator/pkg/service"
)

// CalculatorReconciler reconciles a Calculator object.
type CalculatorReconciler struct {
	client.Client
	Service service.Repository
	Scheme  *runtime.Scheme
}

//+kubebuilder:rbac:groups=calc.example.com,resources=calculators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=calc.example.com,resources=calculators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=calc.example.com,resources=calculators/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CalculatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling Calculator")

	calc, err := r.getCalculator(ctx, req.Name, req.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Info("Got Calculator", "x", calc.Spec.X, "\"y\"", calc.Spec.Y)

	// Process the calculator.
	err = r.Service.ProcessCalculator(ctx, calc)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to process calculator: %w", err)
	}

	// Save the status.
	err = r.manageCalculator(ctx, calc)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create a secret with the result.
	secret, err := r.Service.DefineSecret(ctx, req.Name, req.Namespace, calc.Status.Result)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to define secret: %w", err)
	}

	err = r.manageSecret(ctx, secret)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

//+kubebuilder:rbac:groups=calc.example.com,resources=calculators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=calc.example.com,resources=calculators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=calc.example.com,resources=calculators/finalizers,verbs=update

func (r *CalculatorReconciler) getCalculator(ctx context.Context, name string, namespace string,
) (*calcv1alpha1.Calculator, error) {
	logger := log.FromContext(ctx)

	calc := &calcv1alpha1.Calculator{}

	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, calc)
	if err != nil {
		return nil, fmt.Errorf("failed to get calculator: %w", err)
	}

	logger.Info("Found a calculator", "name", calc.Name, "namespace", calc.Namespace)

	return calc, nil
}

// +kubebuilder:rbac:groups=calc.example.com,resources=calculators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=calc.example.com,resources=calculators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=calc.example.com,resources=calculators/finalizers,verbs=update

func (r *CalculatorReconciler) manageCalculator(ctx context.Context, calc *calcv1alpha1.Calculator) error {
	logger := log.FromContext(ctx)

	logger.Info("Saving the status", "status", calc.Status)

	err := r.Status().Update(ctx, calc)
	if err != nil {
		return fmt.Errorf("failed to update calculator status: %w", err)
	}

	return nil
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *CalculatorReconciler) manageSecret(ctx context.Context, secret *corev1.Secret) error {
	logger := log.FromContext(ctx)

	secretCopy := *secret

	err := r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, &secretCopy)
	if err != nil && errors.IsNotFound(err) {
		logger.Info(fmt.Sprintf("Target secret %s doesn't exist, creating it", secret.Name))

		err = r.Create(ctx, secret)
		if err != nil {
			return fmt.Errorf("failed to create secret: %w", err)
		}

		return nil
	}

	logger.Info(fmt.Sprintf("Target secret %s exists, updating it now", secret.Name))

	err = r.Update(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CalculatorReconciler) SetupWithManager(mgr ctrl.Manager, service service.Repository) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&calcv1alpha1.Calculator{}).
		Complete(r)
	if err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

	r.Service = service

	return nil
}
