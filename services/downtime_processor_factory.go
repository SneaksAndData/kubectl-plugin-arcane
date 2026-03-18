package services

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
)

type DowntimeProcessorFactory struct {
	reader interfaces.UnstructuredReader
}

func NewDowntimeProcessorFactory(reader interfaces.UnstructuredReader) *DowntimeProcessorFactory {
	return &DowntimeProcessorFactory{
		reader: reader,
	}
}

func (s DowntimeProcessorFactory) DowntimeDeclareProcessor(parameters *models.DowntimeDeclareParameters) interfaces.UnstructuredProcessor {
	return &downtimeDeclareProcessor{
		key:         parameters.DowntimeKey,
		reader:      s.reader,
		streamClass: parameters.StreamClass,
	}
}

func (s DowntimeProcessorFactory) DowntimeStopProcessor(parameters *models.DowntimeStopParameters) interfaces.UnstructuredProcessor {
	return &downtimeStopProcessor{
		key:         parameters.DowntimeKey,
		reader:      s.reader,
		streamClass: parameters.StreamClass,
	}
}

func (s DowntimeProcessorFactory) DowntimeSummarizationProcessor(parameters *models.DowntimeListParameters) *DowntimeSummarizationProcessor {
	return &DowntimeSummarizationProcessor{
		reader:      s.reader,
		streamClass: parameters.StreamClass,
	}
}
