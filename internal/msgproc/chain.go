package msgproc

// NewProcessingChain returns chain of the MessageProcessor instances.
func NewProcessingChain(processors ...MessageProcessor) MessageProcessor {
	return &processingChain{
		processors: processors,
	}
}

type processingChain struct {
	processors []MessageProcessor
}

func (p *processingChain) Process(msg string) string {
	if len(p.processors) == 0 {
		return msg
	}

	for _, proc := range p.processors {
		msg = proc.Process(msg)
	}

	return msg
}
