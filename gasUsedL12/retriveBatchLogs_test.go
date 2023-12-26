package gasUsedL12

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrive(t *testing.T) {
	//err := WatchBatches(17505085, 18055390)
	err := WatchBatches(17812792, 17812793)
	require.NoError(t, err)
}

func TestAnalysisSeq(t *testing.T) {
	err := Layer1BatchDataAnalysis(17805237, 17805244)
	require.NoError(t, err)
}

func TestFloat(t *testing.T) {
	var a float64 = 3.612
	var b float64 = 4.2
	fmt.Println(a / b)
}
