package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"zeroscaler.com/core/internal/scaler"
)

type ProxyBackend struct {
	ServiceURL     string
	Deployment     *appsv1.Deployment
	LastActive     time.Time
	ScaleDownAfter metav1.Duration
}

func NewProxy(ctx context.Context, scaler *scaler.Scaler) *Proxy {
	p := &Proxy{
		backends: make(map[string]ProxyBackend, 1),
		scaler:   scaler,
	}
	p.logger = log.FromContext(ctx)

	retryTransport := &RetryTransport{
		Transport: http.DefaultTransport,
		Retries:   50,
		Timeout:   100 * time.Millisecond,
	}

	p.proxy = &httputil.ReverseProxy{
		Transport: retryTransport,
		Director: func(req *http.Request) {
			host := req.Host
			backend, ok := p.backends[host]
			if !ok {
				p.logger.Info("No backend found for host", "host", host)
				return
			}

			// Update the last active time for the backend
			backend.LastActive = time.Now()
			p.backends[host] = backend

			// Check the scaler for the backend's replica count
			replicas, err := p.scaler.GetReplicaCount(backend.Deployment)
			if err != nil {
				p.logger.Error(err, "Failed to get replica count for deployment", "deployment", backend.Deployment.Name)
				return
			}

			// If the replica count is 0, scale the deployment up
			if replicas == 0 {
				p.logger.Info("Scaling up deployment", "deployment", backend.Deployment.Name)
				if err := p.scaler.ScaleDeployment(backend.Deployment, 1); err != nil {
					p.logger.Error(err, "Failed to scale deployment", "deployment", backend.Deployment.Name)
					return
				}
				p.scaler.RefreshDeployment(context.Background(), backend.Deployment)
				p.backends[host] = backend
				// Loop until we have a replica count of 1
				for replicas != 1 {
					// Sleep for 100ms
					time.Sleep(100 * time.Millisecond)
					replicas, err = p.scaler.GetReplicaCount(backend.Deployment)
					if err != nil {
						p.logger.Error(err, "Failed to get replica count for deployment", "deployment", backend.Deployment.Name)
						return
					}
				}
				p.backends[host] = backend
			} else {
				p.logger.Info("Deployment already scaled up", "deployment", backend.Deployment.Name)
			}

			p.logger.Info("Routing request", "host", host, "backend", backend.ServiceURL)
			url, err := url.Parse(backend.ServiceURL)
			if err != nil {
				p.logger.Error(err, "Failed to parse backend URL", "url", backend.ServiceURL)
				return
			}

			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.Host = url.Host
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			p.logger.Error(err, "Error proxying request", "host", r.Host, "path", r.URL.Path)
			// Set the server header on the response
			w.Header().Set("Server", "zeroscaler")
			// Set the response code to 502
			w.WriteHeader(http.StatusBadGateway)
		},
	}
	return p
}

type Proxy struct {
	backends map[string]ProxyBackend
	proxy    *httputil.ReverseProxy
	scaler   *scaler.Scaler
	logger   logr.Logger
}

func (p *Proxy) AddBackend(host, target string, deployment *appsv1.Deployment, scaleDownAfter metav1.Duration) {
	p.logger.Info("Adding backend", "host", host, "target", target)
	p.backends[host] = ProxyBackend{
		ServiceURL:     target,
		Deployment:     deployment,
		LastActive:     time.Now(),
		ScaleDownAfter: scaleDownAfter,
	}
}

func (p *Proxy) RemoveBackend(host string) {
	p.logger.Info("Removing backend", "host", host)
	delete(p.backends, host)
}

func (p *Proxy) Start() error {
	// Kick off a background job that will periodically check the last active time for each backend
	// and scale down the deployment if it has been inactive for a certain period of time
	go func() {
		for {
			for host, backend := range p.backends {
				if time.Since(backend.LastActive) >= backend.ScaleDownAfter.Duration && backend.Deployment.Spec.Replicas != nil && *backend.Deployment.Spec.Replicas > 0 {
					p.logger.Info("Scaling down deployment", "deployment", backend.Deployment.Name, "host", host)
					if err := p.scaler.ScaleDeployment(backend.Deployment, 0); err != nil {
						p.logger.Error(err, "Failed to scale down deployment", "deployment", backend.Deployment.Name)
						p.scaler.RefreshDeployment(context.Background(), backend.Deployment)
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p.logger.Info("Received request", "host", r.Host, "path", r.URL.Path)
		p.proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(":8085", nil)
}

type RetryTransport struct {
	Transport http.RoundTripper
	Retries   int
	Timeout   time.Duration
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < t.Retries; i++ {
		ctx, cancel := context.WithTimeout(req.Context(), t.Timeout)
		defer cancel()

		reqWithTimeout := req.WithContext(ctx)
		resp, err = t.Transport.RoundTrip(reqWithTimeout)
		if err == nil {
			return resp, nil
		}
	}
	return nil, err
}
