//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package osu

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gvallee/go_benchmark/pkg/benchmark"
	"github.com/gvallee/go_hpc_jobmgr/pkg/implem"
	"github.com/gvallee/go_osu/pkg/results"
	"github.com/gvallee/go_software_build/pkg/app"
	"github.com/gvallee/go_software_build/pkg/builder"
	fsutil "github.com/gvallee/go_util/pkg/util"
	"github.com/gvallee/go_workspace/pkg/workspace"
)

const (
	// OSUBaseDir is the name of the directory where the standard OSU is installed
	OSUBaseDir = "OSU"

	// OSUNonConfigMemBaseDir is the name of the directory where the modified OSU for non-contiguous memory is installed
	OSUNonConfigMemBaseDir = "osu_noncontig_mem"

	alltoallBinID   = "alltoall"
	alltoallBinName = "osu_alltoall"

	ialltoallBinID   = "ialltoall"
	ialltoallBinName = "osu_ialltoall"

	alltoallvBinID   = "alltoallv"
	alltoallvBinName = "osu_alltoallv"

	ialltoallvBinID   = "ialltoallv"
	ialltoallvBinName = "osu_ialltoallv"

	alltoallwBinID   = "alltoallw"
	alltoallwBinName = "osu_alltoallw"

	ialltoallwBinID   = "ialltoallw"
	ialltoallwBinName = "osu_ialltoallw"

	allgatherBinID   = "allgather"
	allgatherBinName = "osu_allgather"

	iallgatherBinID   = "iallgather"
	iallgatherBinName = "osu_iallgather"

	allgathervBinID   = "allgatherv"
	allgathervBinName = "osu_allgatherv"

	iallgathervBinID   = "iallgatherv"
	iallgathervBinName = "osu_iallgatherv"

	allreduceBinID   = "allreduce"
	allreduceBinName = "osu_allreduce"

	iallreduceBinID   = "iallreduce"
	iallreduceBinName = "osu_iallreduce"

	barrierBinID   = "barrier"
	barrierBinName = "osu_barrier"

	ibarrierBinID   = "ibarrier"
	ibarrierBinName = "osu_ibarrier"

	bcastBinID   = "bcast"
	bcastBinName = "osu_bcast"

	ibcastBinID   = "ibcast"
	ibcastBinName = "osu_ibcast"

	gatherBinID   = "gather"
	gatherBinName = "osu_gather"

	igatherBinID   = "igather"
	igatherBinName = "osu_igather"

	gathervBinID   = "gatherv"
	gathervBinName = "osu_gatherv"

	igathervBinID   = "igatherv"
	igathervBinName = "osu_igatherv"

	reduceBinID   = "reduce"
	reduceBinName = "osu_reduce"

	ireduceBinID   = "ireduce"
	ireduceBinName = "osu_ireduce"

	reduceScatterBinID   = "reduce_scatter"
	reduceScatterBinName = "osu_reduce_scatter"

	ireduceScatterBinID   = "ireduce_scatter"
	ireduceScatterBinName = "osu_ireduce_scatter"

	scatterBinID   = "scatter"
	scatterBinName = "osu_scatter"

	iscatterBinID   = "iscatter"
	iscatterBinName = "osu_iscatter"

	scattervBinID   = "scatterv"
	scattervBinName = "osu_scatterv"

	iscattervBinID   = "iscatterv"
	iscattervBinName = "osu_iscatterv"

	bwBinID   = "bw"
	bwBinName = "osu_bw"

	latencyBinID   = "latency"
	LatencyBinName = "osu_latency"
)

type SubBenchmark struct {
	AppInfo      app.Info
	OutputParser func(string) (*results.Result, error)
}

// GetPt2PtSubBenchmarks returns the list of all the point-to-point sub-benchmarks from OSU.
// The following arguments must be provided: basedir, the name of the OSU install directory in
// the workspace (should be "osu" by default), the workspace used, and an initialized map where
// the key is the name of the sub-benchmark (e.g., bw) and the value the structure
// describing the benchmark.
func GetPt2PtSubBenchmarks(basedir string, wp *workspace.Config, m map[string]SubBenchmark /*app.Info*/) {
	p2pRelativePath := filepath.Join("libexec", "osu-micro-benchmarks", "mpi", "pt2pt")
	osuPt2Pt2Dir := filepath.Join(wp.InstallDir, basedir, p2pRelativePath)

	bwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    bwBinID,
			BinName: bwBinName,
			BinPath: filepath.Join(osuPt2Pt2Dir, bwBinName),
		},
		OutputParser: nil,
	}
	m[bwBinID] = bwInfo

	latencyInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    latencyBinID,
			BinName: LatencyBinName,
			BinPath: filepath.Join(osuPt2Pt2Dir, LatencyBinName),
		},
		OutputParser: nil,
	}
	m[latencyBinID] = latencyInfo
}

// GetCollectiveSubBenchmarks returns the list of all the collective sub-benchmarks from OSU.
// The following arguments must be provided: basedir, the name of the OSU install directory in
// the workspace (should be "osu" by default), the workspace used, and an initialized map where
// the key is the name of the sub-benchmark (e.g., alltoallv) and the value the structure
// describing the benchmark.
func GetCollectiveSubBenchmarks(basedir string, wp *workspace.Config, m map[string]SubBenchmark) {
	collectiveRelativePath := filepath.Join("libexec", "osu-micro-benchmarks", "mpi", "collective")
	osuCollectiveDir := filepath.Join(wp.InstallDir, basedir, collectiveRelativePath)

	allgatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    allgatherBinID,
			BinName: allgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, allgatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[allgatherBinID] = allgatherInfo

	iallgatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    iallgatherBinID,
			BinName: iallgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallgatherBinName),
		},
		OutputParser: nil,
	}
	m[iallgatherBinID] = iallgatherInfo

	allgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    allgathervBinID,
			BinName: allgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, allgathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[allgathervBinID] = allgathervInfo

	iallgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    iallgathervBinID,
			BinName: iallgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallgathervBinName),
		},
		OutputParser: nil,
	}
	m[iallgathervBinID] = iallgathervInfo

	allreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    allreduceBinID,
			BinName: allreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, allreduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[allreduceBinID] = allreduceInfo

	iallreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    iallreduceBinID,
			BinName: iallreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallreduceBinName),
		},
		OutputParser: nil,
	}
	m[iallreduceBinID] = iallreduceInfo

	alltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    alltoallBinID,
			BinName: alltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[alltoallBinID] = alltoallInfo

	ialltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ialltoallBinID,
			BinName: ialltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallBinName),
		},
		OutputParser: nil,
	}
	m[ialltoallBinID] = ialltoallInfo

	alltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    alltoallvBinID,
			BinName: alltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallvBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[alltoallvBinID] = alltoallvInfo

	ialltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ialltoallvBinID,
			BinName: ialltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallvBinName),
		},
		OutputParser: nil,
	}
	m[ialltoallvBinID] = ialltoallvInfo

	alltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    alltoallwBinID,
			BinName: alltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallwBinName),
		},
		OutputParser: nil,
	}
	m[alltoallwBinID] = alltoallwInfo

	ialltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ialltoallwBinID,
			BinName: ialltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallwBinName),
		},
		OutputParser: nil,
	}
	m[ialltoallwBinID] = ialltoallwInfo

	barrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    barrierBinID,
			BinName: barrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, barrierBinName),
		},
		OutputParser: nil,
	}
	m[barrierBinID] = barrierInfo

	ibarrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ibarrierBinID,
			BinName: ibarrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, ibarrierBinName),
		},
		OutputParser: nil,
	}
	m[ibarrierBinID] = ibarrierInfo

	bcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    bcastBinID,
			BinName: bcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, bcastBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[bcastBinID] = bcastInfo

	ibcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ibcastBinID,
			BinName: ibcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, ibcastBinName),
		},
		OutputParser: nil,
	}
	m[ibcastBinID] = ibcastInfo

	gatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    gatherBinID,
			BinName: gatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, gatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[gatherBinID] = gatherInfo

	igatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    igatherBinID,
			BinName: igatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, igatherBinName),
		},
		OutputParser: nil,
	}
	m[igatherBinID] = igatherInfo

	gathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    gathervBinID,
			BinName: gathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, gathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[gathervBinID] = gathervInfo

	igathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    igathervBinID,
			BinName: igathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, igathervBinName),
		},
		OutputParser: nil,
	}
	m[igathervBinID] = igathervInfo

	reduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    reduceBinID,
			BinName: reduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, reduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[reduceBinID] = reduceInfo

	ireduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ireduceBinID,
			BinName: ireduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, ireduceBinName),
		},
		OutputParser: nil,
	}
	m[ireduceBinID] = ireduceInfo

	reduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    reduceScatterBinID,
			BinName: reduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, reduceScatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[reduceScatterBinID] = reduceScatterInfo

	ireduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ireduceScatterBinID,
			BinName: ireduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, ireduceScatterBinName),
		},
		OutputParser: nil,
	}
	m[ireduceScatterBinID] = ireduceScatterInfo

	scatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    scatterBinID,
			BinName: scatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[scatterBinID] = scatterInfo

	iscatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    iscatterBinID,
			BinName: iscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
		},
		OutputParser: nil,
	}
	m[iscatterBinID] = iscatterInfo

	scattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    scattervBinID,
			BinName: scatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[scattervBinID] = scattervInfo

	iscattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    iscattervBinID,
			BinName: iscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
		},
		OutputParser: nil,
	}
	m[iscattervBinID] = iscattervInfo
}

func getSubBenchmarks(cfg *benchmark.Config, wp *workspace.Config, basedir string) map[string]SubBenchmark {
	m := make(map[string]SubBenchmark)
	GetCollectiveSubBenchmarks(basedir, wp, m)
	GetPt2PtSubBenchmarks(basedir, wp, m)
	return m
}

// ParseCfg is the function to invoke to parse lines from the main configuration files
// that are specific to OSU
func ParseCfg(cfg *benchmark.Config, basedir string, ID string, srcDir string, key string, value string) {
	switch key {
	case "URL":
		cfg.URL = strings.ReplaceAll(value, ID, basedir)
	}

	osuTarballFilename := filepath.Base(cfg.URL)
	cfg.Tarball = filepath.Join(basedir, srcDir, osuTarballFilename)
}

// Compile downloads and installs OSU on the host
func Compile(cfg *benchmark.Config, wp *workspace.Config, flavorName string) (*benchmark.Install, error) {
	b := new(builder.Builder)
	b.Persistent = wp.InstallDir
	b.App.Name = flavorName
	b.App.URL = cfg.URL

	if wp.ScratchDir == "" || wp.InstallDir == "" || wp.BuildDir == "" || wp.SrcDir == "" {
		return nil, fmt.Errorf("invalid workspace")
	}
	b.Env.ScratchDir = wp.ScratchDir
	b.Env.InstallDir = wp.InstallDir
	b.Env.BuildDir = wp.BuildDir
	b.Env.SrcDir = wp.SrcDir

	// fixme: ultimately, we want a persistent install and instead of passing in
	// 'true' to say so, we want to pass in the install directory
	err := b.Load(false)
	if err != nil {
		return nil, err
	}

	// Find MPI and make sure we pass the information about it to the builder
	mpiInfo := new(implem.Info)
	mpiInfo.InstallDir = wp.MpiDir
	err = mpiInfo.Load(nil)
	if err != nil {
		return nil, fmt.Errorf("no suitable MPI available: %s", err)
	}

	log.Printf("- Compiling OSU with MPI from %s", mpiInfo.InstallDir)
	pathEnv := os.Getenv("PATH")
	pathEnv = "PATH=" + filepath.Join(mpiInfo.InstallDir, "bin") + ":" + pathEnv
	b.Env.Env = append(b.Env.Env, pathEnv)
	ldPathEnv := os.Getenv("LD_LIBRARY_PATH")
	ldPathEnv = "LD_LIBRARY_PATH=" + filepath.Join(mpiInfo.InstallDir, "lib") + ":" + ldPathEnv
	b.Env.Env = append(b.Env.Env, ldPathEnv)
	ccEnv := "CC=" + filepath.Join(mpiInfo.InstallDir, "bin", "mpicc")
	b.Env.Env = append(b.Env.Env, ccEnv)
	cxxEnv := "CXX=" + filepath.Join(mpiInfo.InstallDir, "bin", "mpicxx")
	b.Env.Env = append(b.Env.Env, cxxEnv)

	// Finally we can install OSU, it is all ready
	res := b.Install()
	if res.Err != nil {
		log.Printf("error while install OSU: %s; stdout: %s; stderr: %s", res.Err, res.Stdout, res.Stderr)
		return nil, res.Err
	}

	installInfo := new(benchmark.Install)
	installInfo.SubBenchmarks = append(installInfo.SubBenchmarks, b.App)
	return installInfo, nil
}

// DetectInstall scans a specific workspace and detect which OSU benchmarks are
// available. It sets the path to the binary during the detection so it is easier
// to start benchmarks later on.
func DetectInstall(cfg *benchmark.Config, wp *workspace.Config) map[string]*benchmark.Install {
	// We manage different flavors of OSU, e.g., classic OSU and OSU for non-contiguous memory
	basedirs := []string{OSUBaseDir, OSUNonConfigMemBaseDir}
	installedOSUs := make(map[string]*benchmark.Install)

	log.Println("Detecting OSU installation...")
	for _, basedir := range basedirs {
		installInfo := new(benchmark.Install)
		osuBenchmarks := getSubBenchmarks(cfg, wp, basedir)
		for benchmarkName, benchmarkInfo := range osuBenchmarks {
			log.Printf("-> Checking if %s is installed...", benchmarkName)
			if fsutil.FileExists(benchmarkInfo.AppInfo.BinPath) {
				log.Printf("\t%s exists", benchmarkInfo.AppInfo.BinPath)
				installInfo.SubBenchmarks = append(installInfo.SubBenchmarks, benchmarkInfo.AppInfo)
			} else {
				log.Printf("\t%s does not exist", benchmarkInfo.AppInfo.BinPath)
			}
		}
		installedOSUs[basedir] = installInfo
	}
	return installedOSUs
}

// Display shows the current OSU configuration
func Display(cfg *benchmark.Config) {
	fmt.Printf("\tOSU URL: %s\n", cfg.URL)
}
