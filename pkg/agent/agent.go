package agent

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"

	agentclient "github.com/kubeshop/testkube/pkg/agent/client"
	"github.com/kubeshop/testkube/pkg/executor/output"

	"github.com/kubeshop/testkube/internal/config"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/cloud"
	"github.com/kubeshop/testkube/pkg/featureflags"
)

const (
	clusterIDMeta          = "cluster-id"
	cloudMigrateMeta       = "migrate"
	orgIdMeta              = "organization-id"
	envIdMeta              = "environment-id"
	healthcheckCommand     = "healthcheck"
	dockerImageVersionMeta = "docker-image-version"
)

// buffer up to five messages per worker
const bufferSizePerWorker = 5

type Agent struct {
	client  cloud.TestKubeCloudAPIClient
	handler fasthttp.RequestHandler
	logger  *zap.SugaredLogger
	apiKey  string

	workerCount    int
	requestBuffer  chan *cloud.ExecuteRequest
	responseBuffer chan *cloud.ExecuteResponse

	logStreamWorkerCount    int
	logStreamRequestBuffer  chan *cloud.LogsStreamRequest
	logStreamResponseBuffer chan *cloud.LogsStreamResponse
	logStreamFunc           func(ctx context.Context, executionID string) (chan output.Output, error)

	testWorkflowNotificationsWorkerCount    int
	testWorkflowNotificationsRequestBuffer  chan *cloud.TestWorkflowNotificationsRequest
	testWorkflowNotificationsResponseBuffer chan *cloud.TestWorkflowNotificationsResponse
	testWorkflowNotificationsFunc           func(ctx context.Context, executionID string) (<-chan testkube.TestWorkflowExecutionNotification, error)

	testWorkflowServiceNotificationsWorkerCount    int
	testWorkflowServiceNotificationsRequestBuffer  chan *cloud.TestWorkflowServiceNotificationsRequest
	testWorkflowServiceNotificationsResponseBuffer chan *cloud.TestWorkflowServiceNotificationsResponse
	testWorkflowServiceNotificationsFunc           func(ctx context.Context, executionID, serviceName string, serviceIndex int) (<-chan testkube.TestWorkflowExecutionNotification, error)

	testWorkflowParallelStepNotificationsWorkerCount    int
	testWorkflowParallelStepNotificationsRequestBuffer  chan *cloud.TestWorkflowParallelStepNotificationsRequest
	testWorkflowParallelStepNotificationsResponseBuffer chan *cloud.TestWorkflowParallelStepNotificationsResponse
	testWorkflowParallelStepNotificationsFunc           func(ctx context.Context, executionID, ref string, workerIndex int) (<-chan testkube.TestWorkflowExecutionNotification, error)

	events              chan testkube.Event
	sendTimeout         time.Duration
	receiveTimeout      time.Duration
	healthcheckInterval time.Duration

	clusterID          string
	clusterName        string
	features           featureflags.FeatureFlags
	dockerImageVersion string

	proContext *config.ProContext
}

func NewAgent(logger *zap.SugaredLogger,
	handler fasthttp.RequestHandler,
	client cloud.TestKubeCloudAPIClient,
	logStreamFunc func(ctx context.Context, executionID string) (chan output.Output, error),
	workflowNotificationsFunc func(ctx context.Context, executionID string) (<-chan testkube.TestWorkflowExecutionNotification, error),
	workflowServiceNotificationsFunc func(ctx context.Context, executionID, serviceName string, serviceIndex int) (<-chan testkube.TestWorkflowExecutionNotification, error),
	workflowParallelStepNotificationsFunc func(ctx context.Context, executionID, ref string, workerIndex int) (<-chan testkube.TestWorkflowExecutionNotification, error),
	clusterID string,
	clusterName string,
	features featureflags.FeatureFlags,
	proContext *config.ProContext,
	dockerImageVersion string,
) (*Agent, error) {
	return &Agent{
		handler:                                 handler,
		logger:                                  logger.With("service", "Agent", "environmentId", proContext.EnvID),
		apiKey:                                  proContext.APIKey,
		client:                                  client,
		events:                                  make(chan testkube.Event),
		workerCount:                             proContext.WorkerCount,
		requestBuffer:                           make(chan *cloud.ExecuteRequest, bufferSizePerWorker*proContext.WorkerCount),
		responseBuffer:                          make(chan *cloud.ExecuteResponse, bufferSizePerWorker*proContext.WorkerCount),
		receiveTimeout:                          5 * time.Minute,
		sendTimeout:                             30 * time.Second,
		healthcheckInterval:                     30 * time.Second,
		logStreamWorkerCount:                    proContext.LogStreamWorkerCount,
		logStreamRequestBuffer:                  make(chan *cloud.LogsStreamRequest, bufferSizePerWorker*proContext.LogStreamWorkerCount),
		logStreamResponseBuffer:                 make(chan *cloud.LogsStreamResponse, bufferSizePerWorker*proContext.LogStreamWorkerCount),
		logStreamFunc:                           logStreamFunc,
		testWorkflowNotificationsWorkerCount:    proContext.WorkflowNotificationsWorkerCount,
		testWorkflowNotificationsRequestBuffer:  make(chan *cloud.TestWorkflowNotificationsRequest, bufferSizePerWorker*proContext.WorkflowNotificationsWorkerCount),
		testWorkflowNotificationsResponseBuffer: make(chan *cloud.TestWorkflowNotificationsResponse, bufferSizePerWorker*proContext.WorkflowNotificationsWorkerCount),
		testWorkflowNotificationsFunc:           workflowNotificationsFunc,
		testWorkflowServiceNotificationsWorkerCount:         proContext.WorkflowServiceNotificationsWorkerCount,
		testWorkflowServiceNotificationsRequestBuffer:       make(chan *cloud.TestWorkflowServiceNotificationsRequest, bufferSizePerWorker*proContext.WorkflowServiceNotificationsWorkerCount),
		testWorkflowServiceNotificationsResponseBuffer:      make(chan *cloud.TestWorkflowServiceNotificationsResponse, bufferSizePerWorker*proContext.WorkflowServiceNotificationsWorkerCount),
		testWorkflowServiceNotificationsFunc:                workflowServiceNotificationsFunc,
		testWorkflowParallelStepNotificationsWorkerCount:    proContext.WorkflowParallelStepNotificationsWorkerCount,
		testWorkflowParallelStepNotificationsRequestBuffer:  make(chan *cloud.TestWorkflowParallelStepNotificationsRequest, bufferSizePerWorker*proContext.WorkflowParallelStepNotificationsWorkerCount),
		testWorkflowParallelStepNotificationsResponseBuffer: make(chan *cloud.TestWorkflowParallelStepNotificationsResponse, bufferSizePerWorker*proContext.WorkflowParallelStepNotificationsWorkerCount),
		testWorkflowParallelStepNotificationsFunc:           workflowParallelStepNotificationsFunc,
		clusterID:          clusterID,
		clusterName:        clusterName,
		features:           features,
		proContext:         proContext,
		dockerImageVersion: dockerImageVersion,
	}, nil
}

func (ag *Agent) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := ag.run(ctx)

		ag.logger.Errorw("agent connection failed, reconnecting", "error", err)

		// TODO: some smart back off strategy?
		time.Sleep(5 * time.Second)
	}
}

func (ag *Agent) run(ctx context.Context) (err error) {
	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ag.runCommandLoop(groupCtx)
	})

	g.Go(func() error {
		return ag.runWorkers(groupCtx, ag.workerCount)
	})

	g.Go(func() error {
		return ag.runEventLoop(groupCtx)
	})

	if !ag.features.LogsV2 {
		g.Go(func() error {
			return ag.runLogStreamLoop(groupCtx)
		})
		g.Go(func() error {
			return ag.runLogStreamWorker(groupCtx, ag.logStreamWorkerCount)
		})
	}

	g.Go(func() error {
		return ag.runTestWorkflowNotificationsLoop(groupCtx)
	})
	g.Go(func() error {
		return ag.runTestWorkflowNotificationsWorker(groupCtx, ag.testWorkflowNotificationsWorkerCount)
	})

	g.Go(func() error {
		return ag.runTestWorkflowServiceNotificationsLoop(groupCtx)
	})
	g.Go(func() error {
		return ag.runTestWorkflowServiceNotificationsWorker(groupCtx, ag.testWorkflowServiceNotificationsWorkerCount)
	})

	g.Go(func() error {
		return ag.runTestWorkflowParallelStepNotificationsLoop(groupCtx)
	})
	g.Go(func() error {
		return ag.runTestWorkflowParallelStepNotificationsWorker(groupCtx, ag.testWorkflowParallelStepNotificationsWorkerCount)
	})

	err = g.Wait()

	return err
}

func (ag *Agent) sendResponse(ctx context.Context, stream cloud.TestKubeCloudAPI_ExecuteClient, resp *cloud.ExecuteResponse) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- stream.Send(resp)
		close(errChan)
	}()

	t := time.NewTimer(ag.sendTimeout)
	select {
	case err := <-errChan:
		if !t.Stop() {
			<-t.C
		}
		return err
	case <-ctx.Done():
		if !t.Stop() {
			<-t.C
		}

		return ctx.Err()
	case <-t.C:
		return errors.New("send response too slow")
	}
}

func (ag *Agent) receiveCommand(ctx context.Context, stream cloud.TestKubeCloudAPI_ExecuteClient) (*cloud.ExecuteRequest, error) {
	respChan := make(chan cloudResponse, 1)
	go func() {
		cmd, err := stream.Recv()
		respChan <- cloudResponse{resp: cmd, err: err}
	}()

	t := time.NewTimer(ag.receiveTimeout)
	var cmd *cloud.ExecuteRequest
	select {
	case resp := <-respChan:
		if !t.Stop() {
			<-t.C
		}

		cmd = resp.resp
		err := resp.err

		if err != nil {
			ag.logger.Errorf("agent stream receive: %v", err)
			return nil, err
		}
	case <-ctx.Done():
		if !t.Stop() {
			<-t.C
		}

		return nil, ctx.Err()
	case <-t.C:
		return nil, errors.New("stream receive too slow")
	}

	return cmd, nil
}

func (ag *Agent) runCommandLoop(ctx context.Context) error {
	if ag.proContext.APIKey != "" {
		ctx = agentclient.AddAPIKeyMeta(ctx, ag.proContext.APIKey)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, clusterIDMeta, ag.clusterID)
	ctx = metadata.AppendToOutgoingContext(ctx, cloudMigrateMeta, ag.proContext.Migrate)
	ctx = metadata.AppendToOutgoingContext(ctx, envIdMeta, ag.proContext.EnvID)
	ctx = metadata.AppendToOutgoingContext(ctx, orgIdMeta, ag.proContext.OrgID)
	ctx = metadata.AppendToOutgoingContext(ctx, dockerImageVersionMeta, ag.dockerImageVersion)

	ag.logger.Infow("initiating streaming connection with control plane")
	// creates a new Stream from the client side. ctx is used for the lifetime of the stream.
	opts := []grpc.CallOption{grpc.UseCompressor(gzip.Name), grpc.MaxCallRecvMsgSize(math.MaxInt32)}
	stream, err := ag.client.ExecuteAsync(ctx, opts...)
	if err != nil {
		ag.logger.Errorf("failed to execute: %w", err)
		return errors.Wrap(err, "failed to setup stream")
	}

	// GRPC stream have special requirements for concurrency on SendMsg, and RecvMsg calls.
	// Please check https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md
	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for {
			cmd, err := ag.receiveCommand(groupCtx, stream)
			if err != nil {
				return err
			}

			ag.requestBuffer <- cmd
		}
	})

	g.Go(func() error {
		for {
			select {
			case resp := <-ag.responseBuffer:
				err := ag.sendResponse(groupCtx, stream, resp)
				if err != nil {
					return err
				}
			case <-groupCtx.Done():
				return groupCtx.Err()
			}
		}
	})

	err = g.Wait()

	return err
}

func (ag *Agent) runWorkers(ctx context.Context, numWorkers int) error {
	g, groupCtx := errgroup.WithContext(ctx)
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case cmd := <-ag.requestBuffer:
					select {
					case ag.responseBuffer <- ag.executeCommand(groupCtx, cmd):
					case <-groupCtx.Done():
						return groupCtx.Err()
					}
				case <-groupCtx.Done():
					return groupCtx.Err()
				}
			}
		})
	}
	return g.Wait()
}

func (ag *Agent) executeCommand(_ context.Context, cmd *cloud.ExecuteRequest) *cloud.ExecuteResponse {
	switch {
	case cmd.Url == healthcheckCommand:
		return &cloud.ExecuteResponse{MessageId: cmd.MessageId, Status: 0}
	default:
		req := &fasthttp.RequestCtx{}
		r := fasthttp.AcquireRequest()
		r.Header.SetHost("localhost")
		r.Header.SetMethod(cmd.Method)

		for k, values := range cmd.Headers {
			for _, value := range values.Header {
				r.Header.Add(k, value)
			}
		}
		r.SetBody(cmd.Body)
		uri := &fasthttp.URI{}

		err := uri.Parse(nil, []byte(cmd.Url))
		if err != nil {
			ag.logger.Errorf("agent bad command url: %w", err)
			resp := &cloud.ExecuteResponse{MessageId: cmd.MessageId, Status: 400, Body: []byte(fmt.Sprintf("bad command url: %s", err))}
			return resp
		}
		r.SetURI(uri)

		req.Init(r, nil, nil)
		ag.handler(req)

		fasthttp.ReleaseRequest(r)

		headers := make(map[string]*cloud.HeaderValue)
		req.Response.Header.VisitAll(func(key, value []byte) {
			_, ok := headers[string(key)]
			if !ok {
				headers[string(key)] = &cloud.HeaderValue{Header: []string{string(value)}}
				return
			}

			headers[string(key)].Header = append(headers[string(key)].Header, string(value))
		})

		resp := &cloud.ExecuteResponse{MessageId: cmd.MessageId, Headers: headers, Status: int64(req.Response.StatusCode()), Body: req.Response.Body()}

		return resp
	}
}

type cloudResponse struct {
	resp *cloud.ExecuteRequest
	err  error
}
