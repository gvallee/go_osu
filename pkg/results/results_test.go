package results

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	alltoallResultsExample = `/global/scratch/users/geoffroy/a2av_validation_ws/install/ompi/bin/mpirun

# OSU MPI All-to-Allv Personalized Exchange Latency Test v5.6.3
# Size       Avg Latency(us)
1                     159.35
2                     118.59
4                     137.19
8                     156.69
16                    260.03
32                    374.01
64                    640.70
128                  1034.80
256                  1976.88
512                  3359.99
1024                 5056.74
2048                 9060.91
4096                17582.48
8192                34491.39
16384               65969.07
32768              131362.47
65536              264175.86
131072             442734.73
262144             722343.56
524288            1432346.21
`
)

func TestExtractDataFromOutput(t *testing.T) {
	dataSize, values, err := ExtractDataFromOutput(strings.Split(alltoallResultsExample, "\n"))
	if err != nil {
		t.Fatalf("ExtractDataFromOutput() failed: %s", err)
	}
	if len(dataSize) != len(values) {
		t.Fatalf("inconsistent data: %d vs. %d", len(dataSize), len(values))
	}
	expectedSize := 1.0
	for _, sz := range dataSize {
		if sz != expectedSize {
			t.Fatalf("unexpected size: %f when %f was expectede", sz, expectedSize)
		}
		expectedSize *= 2
	}
}

func TestExcelize(t *testing.T) {
	dataSizes, values, err := ExtractDataFromOutput(strings.Split(alltoallResultsExample, "\n"))
	if err != nil {
		t.Fatalf("ExtractDataFromOutput() failed: %s", err)
	}
	res := new(Results)
	err = rawDataToResults(dataSizes, values, res)
	if err != nil {
		t.Fatalf("rawDataToResults() failed: %s", err)
	}

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, "test.xslx")
	err = Excelize(tempFile, res)
	if err != nil {
		t.Fatalf("Excelize() failed: %s", err)
	}
}
