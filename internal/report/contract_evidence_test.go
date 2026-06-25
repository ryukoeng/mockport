package report

import (
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestContractEvidenceIsAdapterAlias(t *testing.T) {
	var _ ContractEvidence = adapter.ContractEvidence{}
	var _ adapter.ContractEvidence = ContractEvidence{}
}
