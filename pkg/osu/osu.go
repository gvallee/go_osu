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

	AlltoallID      = "alltoall"
	alltoallBinName = "osu_alltoall"

	IalltoallID      = "ialltoall"
	ialltoallBinName = "osu_ialltoall"

	AlltoallvID      = "alltoallv"
	alltoallvBinName = "osu_alltoallv"

	IalltoallvID      = "ialltoallv"
	ialltoallvBinName = "osu_ialltoallv"

	AlltoallwID      = "alltoallw"
	alltoallwBinName = "osu_alltoallw"

	IalltoallwID      = "ialltoallw"
	ialltoallwBinName = "osu_ialltoallw"

	AllgatherID      = "allgather"
	allgatherBinName = "osu_allgather"

	IallgatherID      = "iallgather"
	iallgatherBinName = "osu_iallgather"

	AllgathervID      = "allgatherv"
	allgathervBinName = "osu_allgatherv"

	IallgathervID      = "iallgatherv"
	iallgathervBinName = "osu_iallgatherv"

	AllreduceID      = "allreduce"
	allreduceBinName = "osu_allreduce"

	IallreduceID      = "iallreduce"
	iallreduceBinName = "osu_iallreduce"

	BarrierID      = "barrier"
	barrierBinName = "osu_barrier"

	IbarrierID      = "ibarrier"
	ibarrierBinName = "osu_ibarrier"

	BcastID      = "bcast"
	bcastBinName = "osu_bcast"

	IbcastID      = "ibcast"
	ibcastBinName = "osu_ibcast"

	GatherID      = "gather"
	gatherBinName = "osu_gather"

	IgatherID      = "igather"
	igatherBinName = "osu_igather"

	GathervID      = "gatherv"
	gathervBinName = "osu_gatherv"

	IgathervID      = "igatherv"
	igathervBinName = "osu_igatherv"

	ReduceID      = "reduce"
	reduceBinName = "osu_reduce"

	IreduceID      = "ireduce"
	ireduceBinName = "osu_ireduce"

	ReduceScatterID      = "reduce_scatter"
	reduceScatterBinName = "osu_reduce_scatter"

	IreduceScatterID      = "ireduce_scatter"
	ireduceScatterBinName = "osu_ireduce_scatter"

	ScatterID      = "scatter"
	scatterBinName = "osu_scatter"

	IscatterID      = "iscatter"
	iscatterBinName = "osu_iscatter"

	ScattervID      = "scatterv"
	scattervBinName = "osu_scatterv"

	IscattervID      = "iscatterv"
	iscattervBinName = "osu_iscatterv"

	BWID      = "bw"
	bwBinName = "osu_bw"

	LatencyID      = "latency"
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
			Name:    BWID,
			BinName: bwBinName,
			BinPath: filepath.Join(osuPt2Pt2Dir, bwBinName),
		},
		OutputParser: nil,
	}
	m[BWID] = bwInfo

	latencyInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    LatencyID,
			BinName: LatencyBinName,
			BinPath: filepath.Join(osuPt2Pt2Dir, LatencyBinName),
		},
		OutputParser: nil,
	}
	m[LatencyID] = latencyInfo
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
			Name:    AllgatherID,
			BinName: allgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, allgatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllgatherID] = allgatherInfo

	iallgatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallgatherID,
			BinName: iallgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallgatherBinName),
		},
		OutputParser: nil,
	}
	m[IallgatherID] = iallgatherInfo

	allgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AllgathervID,
			BinName: allgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, allgathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllgathervID] = allgathervInfo

	iallgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallgathervID,
			BinName: iallgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallgathervBinName),
		},
		OutputParser: nil,
	}
	m[IallgathervID] = iallgathervInfo

	allreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AllreduceID,
			BinName: allreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, allreduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllreduceID] = allreduceInfo

	iallreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallreduceID,
			BinName: iallreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, iallreduceBinName),
		},
		OutputParser: nil,
	}
	m[IallreduceID] = iallreduceInfo

	alltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallID,
			BinName: alltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AlltoallID] = alltoallInfo

	ialltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallID,
			BinName: ialltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallID] = ialltoallInfo

	alltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallvID,
			BinName: alltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallvBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AlltoallvID] = alltoallvInfo

	ialltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallvID,
			BinName: ialltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallvBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallvID] = ialltoallvInfo

	alltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallwID,
			BinName: alltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, alltoallwBinName),
		},
		OutputParser: nil,
	}
	m[AlltoallwID] = alltoallwInfo

	ialltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallwID,
			BinName: ialltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, ialltoallwBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallwID] = ialltoallwInfo

	barrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    BarrierID,
			BinName: barrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, barrierBinName),
		},
		OutputParser: nil,
	}
	m[BarrierID] = barrierInfo

	ibarrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IbarrierID,
			BinName: ibarrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, ibarrierBinName),
		},
		OutputParser: nil,
	}
	m[IbarrierID] = ibarrierInfo

	bcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    BcastID,
			BinName: bcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, bcastBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[BcastID] = bcastInfo

	ibcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IbcastID,
			BinName: ibcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, ibcastBinName),
		},
		OutputParser: nil,
	}
	m[IbcastID] = ibcastInfo

	gatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    GatherID,
			BinName: gatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, gatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[GatherID] = gatherInfo

	igatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IgatherID,
			BinName: igatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, igatherBinName),
		},
		OutputParser: nil,
	}
	m[IgatherID] = igatherInfo

	gathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    GathervID,
			BinName: gathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, gathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[GathervID] = gathervInfo

	igathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IgathervID,
			BinName: igathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, igathervBinName),
		},
		OutputParser: nil,
	}
	m[IgathervID] = igathervInfo

	reduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ReduceID,
			BinName: reduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, reduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ReduceID] = reduceInfo

	ireduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IreduceID,
			BinName: ireduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, ireduceBinName),
		},
		OutputParser: nil,
	}
	m[IreduceID] = ireduceInfo

	reduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ReduceScatterID,
			BinName: reduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, reduceScatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ReduceScatterID] = reduceScatterInfo

	ireduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IreduceScatterID,
			BinName: ireduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, ireduceScatterBinName),
		},
		OutputParser: nil,
	}
	m[IreduceScatterID] = ireduceScatterInfo

	scatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ScatterID,
			BinName: scatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ScatterID] = scatterInfo

	iscatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IscatterID,
			BinName: iscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
		},
		OutputParser: nil,
	}
	m[IscatterID] = iscatterInfo

	scattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ScattervID,
			BinName: scatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ScattervID] = scattervInfo

	iscattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IscattervID,
			BinName: iscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
		},
		OutputParser: nil,
	}
	m[IscattervID] = iscattervInfo
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

	if wp.ScratchDir == "" {
		return nil, fmt.Errorf("workspace's scratch directory is undefined")
	}

	if wp.InstallDir == "" {
		return nil, fmt.Errorf("workspace's install directory is undefined")
	}

	if wp.BuildDir == "" {
		return nil, fmt.Errorf("workspace's build directory is undefined")
	}

	if wp.SrcDir == "" {
		return nil, fmt.Errorf("workspace's source directory is undefined")
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
		return nil, fmt.Errorf("builder.Install() failed: %w", res.Err)
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
