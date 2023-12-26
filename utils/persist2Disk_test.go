package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPersist2Disk(t *testing.T) {
	err := Save2File("./testwrite.txt", []string{"0x22a45f235889df3f2f2661dc28054160170225609b5b746252b8c3faac63fe89"})
	require.NoError(t, err)
	err = Save2File("./testwrite.txt", []string{"0x22a45f235889df3f2f2661dc28054160170225609b5b746252b8c3faac63fe89"})
	require.NoError(t, err)
	err = Save2File("./testwrite.txt", []string{"21312312312312"})
	require.NoError(t, err)
}
