package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnonimizerJSON(t *testing.T) {
	inputPath := "../testdata/input.json"
	outputPath := "../testdata/output.json"

	inputFile, errOpen := os.Open(inputPath)
	require.NoError(t,
		errOpen,
		"Failed to open input file",
	)
	defer inputFile.Close()

	outputFile, errCr := os.Create(outputPath)
	require.NoError(t,
		errCr,
		"Failed to create output file",
	)
	outputFile.Close()

	require.NoError(t,
		AnonymizeJSON(inputFile, outputFile),
		"AnonymizeJSON failed",
	)
}
