package scaler

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The scaler is responsible for sending requests to the Kubernetes API to scale the target deployment up or down
// based on the number of requests received for the given hosts.
// The scaler also maintains routing tables to route requests to the correct deployment based on the host header.

func NewScaler(client client.Client) *Scaler {
	return &Scaler{
		Client: client,
	}
}

type Scaler struct {
	Client client.Client
}

func (s *Scaler) ScaleDeployment(deployment *appsv1.Deployment, replicas int32) error {
	deployment.Spec.Replicas = &replicas
	return s.Client.Update(context.Background(), deployment)
}

func (s *Scaler) GetReplicaCount(deployment *appsv1.Deployment) (int32, error) {
	err := s.RefreshDeployment(context.Background(), deployment)
	if err != nil {
		return 0, err
	}

	return deployment.Status.ReadyReplicas, nil
}

func (s *Scaler) RefreshDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	return s.Client.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, deployment)
}
