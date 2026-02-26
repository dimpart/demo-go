package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
)

type ReceiptCommandProcessor struct {
	*BaseCommandProcessor
}

// Override
func (cpu *ReceiptCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	// no need to response receipt command
	return nil
}
