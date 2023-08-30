// The http package is responsible for initializing the server, the router with handlers, and for processing requests.
package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"segmentation-service/internal/ports"
	"segmentation-service/pkg/infra/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Adapter struct {
	s          *http.Server
	l          net.Listener
	segmentSvc ports.SegmentService
}

type AdapterOptions struct {
	HTTP_port   int
	Timeout     time.Duration
	IdleTimeout time.Duration
}

// New instantiates the adapter.
func New(segmentService ports.SegmentService, opts AdapterOptions) (*Adapter, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.HTTP_port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}

	router := gin.Default()
	server := http.Server{
		Handler:      router,
		ReadTimeout:  opts.Timeout,
		WriteTimeout: opts.Timeout,
		IdleTimeout:  opts.IdleTimeout, // client connection lifetime
	}
	a := Adapter{
		s:          &server,
		l:          l,
		segmentSvc: segmentService,
	}
	err = initRouter(&a, router)
	return &a, err
}

// Start starts an http server that accepts incoming connections on the Listener.
func (a *Adapter) Start() error {
	logger.Get().Info("starting http server...")

	eg := &errgroup.Group{}
	eg.Go(func() error {
		return a.s.Serve(a.l)
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Stop stops the http server.
func (a *Adapter) Stop(ctx context.Context) error {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		err = a.s.Shutdown(ctx)
	})
	return err
}
