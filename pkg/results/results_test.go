package results

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	alltoallResultsExample1 = `/global/scratch/users/geoffroy/a2av_validation_ws/install/ompi/bin/mpirun

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

	alltoallResultsExample2 = `/global/scratch/users/geoffroy/a2av_validation_ws/install/ompi/bin/mpirun

# OSU MPI All-to-Allv Personalized Exchange Latency Test v5.6.3
# Size       Avg Latency(us)
1                     1
2                     1.2
4                     1.3
8                     1.4
16                    1.5
32                    1.6
64                    1.7
128                  1.8
256                  1.9
512                  2
1024                 2.1
2048                 2.2
4096                2.3
8192                2.4
16384               2.5
32768              2.6
65536              2.7
131072             2.8
262144             2.9
524288            3
`
)

func TestExtractDataFromOutput(t *testing.T) {
	dataSize, values, err := ExtractDataFromOutput(strings.Split(alltoallResultsExample1, "\n"))
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

func prepExcelizeTest(t *testing.T) *Results {
	dataSizes, values, err := ExtractDataFromOutput(strings.Split(alltoallResultsExample1, "\n"))
	if err != nil {
		t.Fatalf("ExtractDataFromOutput() failed: %s", err)
	}
	dataSizes2, values2, err := ExtractDataFromOutput(strings.Split(alltoallResultsExample2, "\n"))
	if err != nil {
		t.Fatalf("ExtractDataFromOutput() faileed: %s", err)
	}

	if len(dataSizes2) != len(dataSizes) {
		t.Fatalf("data sizes of different length: %d vs. %d", len(dataSizes2), len(dataSizes))
	}
	for idx := range dataSizes {
		if dataSizes[idx] != dataSizes2[idx] {
			t.Fatalf("invalid data: %f vs. %f", dataSizes[idx], dataSizes2[idx])
		}
	}

	res := new(Results)
	err = rawDataToResults(dataSizes, values, res)
	if err != nil {
		t.Fatalf("rawDataToResults() failed: %s", err)
	}

	err = rawDataToResults(dataSizes2, values2, res)
	if err != nil {
		t.Fatalf("rawDataToResults() failed: %s", err)
	}

	return res
}

func TestExcelize(t *testing.T) {
	res := prepExcelizeTest(t)

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, "test.xlsx")
	err = Excelize(tempFile, res)
	if err != nil {
		t.Fatalf("Excelize() failed: %s", err)
	}
}

func TestExcelizeWithLabels(t *testing.T) {
	labels := []string{"label1", "label2"}
	res := prepExcelizeTest(t)

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, "test_with_labels.xlsx")
	err = ExcelizeWithLabels(tempFile, res, labels)
	if err != nil {
		t.Fatalf("Excelize() failed: %s", err)
	}
}
