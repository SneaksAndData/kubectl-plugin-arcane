package services

import (
	"context"

	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/filter"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/publisher"
)

// Ensure downtime implements cmdinterfaces.DowntimeService
var _ cmdinterfaces.DowntimeService = (*downtime)(nil)

// downtime is a service that provides downtime operations.
type downtime struct {
	clientProvider cmdinterfaces.ClientProvider
	factory        *DowntimeProcessorFactory
	executionQueue interfaces.ExecutionQueue
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService(clientProvider cmdinterfaces.ClientProvider, factory *DowntimeProcessorFactory) cmdinterfaces.DowntimeService {
	return &downtime{
		clientProvider: clientProvider,
		factory:        factory,
		executionQueue: NewExecutionQueue(clientProvider),
	}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(ctx context.Context, parameters *models.DowntimeDeclareParameters) error {
	f := filter.NewUnsuspendedByNamePrefix(parameters.Prefix)
	membersPublisher := publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, parameters.Namespace, f)
	return s.executionQueue.ProcessQueue(ctx, s.factory.DowntimeDeclareProcessor(parameters), logging.Printer("suspended"), membersPublisher)
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(ctx context.Context, parameters *models.DowntimeStopParameters) error {
	f := filter.NewByDowntimeKey(parameters.DowntimeKey)
	membersPublisher := publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, "", f)
	return s.executionQueue.ProcessQueue(ctx, s.factory.DowntimeStopProcessor(parameters), logging.Printer("started"), membersPublisher)
}

func (s *downtime) GetSummary(ctx context.Context, parameters *models.DowntimeSummaryParameters) (cmdinterfaces.DowntimeSummary, error) {
	var queuePublisher interfaces.QueuePublisher
	if parameters.StreamClass == "" {
		queuePublisher = publisher.NewAllStreamDefinitionsPublisher(s.clientProvider)
	} else {
		queuePublisher = publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, "", filter.NewAllowAll())
	}

	processor := s.factory.DowntimeSummarizationProcessor()
	err := s.executionQueue.ProcessQueue(ctx, processor, logging.Printer(""), queuePublisher)
	if err != nil { // coverage-ignore
		return nil, err
	}

	return NewDowntimeSummary(processor.Summary, processor.Durations), nil
}
