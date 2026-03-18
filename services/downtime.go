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
	ExecutionQueue interfaces.ExecutionQueue
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService(clientProvider cmdinterfaces.ClientProvider, factory *DowntimeProcessorFactory) cmdinterfaces.DowntimeService {
	return &downtime{
		clientProvider: clientProvider,
		factory:        factory,
		ExecutionQueue: &DefinitionProcessor{},
	}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(ctx context.Context, parameters *models.DowntimeDeclareParameters) error {
	f := filter.NewUnsuspendedByNamePrefix(parameters.Prefix)
	membersPublisher := publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, parameters.Namespace, f)
	return s.ExecutionQueue.ProcessQueue(ctx, s.factory.DowntimeDeclareProcessor(parameters), logging.Printer("suspended"), membersPublisher)
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(ctx context.Context, parameters *models.DowntimeStopParameters) error {
	f := filter.NewByDowntimeKey(parameters.DowntimeKey)
	membersPublisher := publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, "", f)
	return s.ExecutionQueue.ProcessQueue(ctx, s.factory.DowntimeStopProcessor(parameters), logging.Printer("started"), membersPublisher)
}

func (s *downtime) ListDowntimes(ctx context.Context, parameters *models.DowntimeListParameters) error {
	var queuePublisher interfaces.QueuePublisher
	if parameters.StreamClass == "" {
		queuePublisher = publisher.NewAllStreamDefinitionsPublisher(s.clientProvider)
	} else {
		queuePublisher = publisher.NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, "", filter.NewAllowAll())
	}

	return s.ExecutionQueue.ProcessQueue(ctx, s.factory.DowntimeSummarizationProcessor(parameters), logging.Printer("started"), queuePublisher)
}
