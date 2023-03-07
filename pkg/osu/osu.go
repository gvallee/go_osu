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
	AlltoallBinName = "osu_alltoall"

	IalltoallID      = "ialltoall"
	IalltoallBinName = "osu_ialltoall"

	AlltoallvID      = "alltoallv"
	AlltoallvBinName = "osu_alltoallv"

	IalltoallvID      = "ialltoallv"
	IalltoallvBinName = "osu_ialltoallv"

	AlltoallwID      = "alltoallw"
	AlltoallwBinName = "osu_alltoallw"

	IalltoallwID      = "ialltoallw"
	IalltoallwBinName = "osu_ialltoallw"

	AllgatherID      = "allgather"
	AllgatherBinName = "osu_allgather"

	IallgatherID      = "iallgather"
	IallgatherBinName = "osu_iallgather"

	AllgathervID      = "allgatherv"
	AllgathervBinName = "osu_allgatherv"

	IallgathervID      = "iallgatherv"
	IallgathervBinName = "osu_iallgatherv"

	AllreduceID      = "allreduce"
	AllreduceBinName = "osu_allreduce"

	IallreduceID      = "iallreduce"
	IallreduceBinName = "osu_iallreduce"

	BarrierID      = "barrier"
	BarrierBinName = "osu_barrier"

	IbarrierID      = "ibarrier"
	IbarrierBinName = "osu_ibarrier"

	BcastID      = "bcast"
	BcastBinName = "osu_bcast"

	IbcastID      = "ibcast"
	IbcastBinName = "osu_ibcast"

	GatherID      = "gather"
	GatherBinName = "osu_gather"

	IgatherID      = "igather"
	IgatherBinName = "osu_igather"

	GathervID      = "gatherv"
	GathervBinName = "osu_gatherv"

	IgathervID      = "igatherv"
	IgathervBinName = "osu_igatherv"

	ReduceID      = "reduce"
	ReduceBinName = "osu_reduce"

	IreduceID      = "ireduce"
	IreduceBinName = "osu_ireduce"

	ReduceScatterID      = "reduce_scatter"
	ReduceScatterBinName = "osu_reduce_scatter"

	IreduceScatterID      = "ireduce_scatter"
	IreduceScatterBinName = "osu_ireduce_scatter"

	ScatterID      = "scatter"
	ScatterBinName = "osu_scatter"

	IscatterID      = "iscatter"
	IscatterBinName = "osu_iscatter"

	ScattervID      = "scatterv"
	ScattervBinName = "osu_scatterv"

	IscattervID      = "iscatterv"
	IscattervBinName = "osu_iscatterv"

	BWID      = "bw"
	BWBinName = "osu_bw"

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
			BinName: BWBinName,
			BinPath: filepath.Join(osuPt2Pt2Dir, BWBinName),
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
			BinName: AllgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, AllgatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllgatherID] = allgatherInfo

	iallgatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallgatherID,
			BinName: IallgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, IallgatherBinName),
		},
		OutputParser: nil,
	}
	m[IallgatherID] = iallgatherInfo

	allgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AllgathervID,
			BinName: AllgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, AllgathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllgathervID] = allgathervInfo

	iallgathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallgathervID,
			BinName: IallgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, IallgathervBinName),
		},
		OutputParser: nil,
	}
	m[IallgathervID] = iallgathervInfo

	allreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AllreduceID,
			BinName: AllreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, AllreduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AllreduceID] = allreduceInfo

	iallreduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IallreduceID,
			BinName: IallreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, IallreduceBinName),
		},
		OutputParser: nil,
	}
	m[IallreduceID] = iallreduceInfo

	alltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallID,
			BinName: AlltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, AlltoallBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AlltoallID] = alltoallInfo

	ialltoallInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallID,
			BinName: IalltoallBinName,
			BinPath: filepath.Join(osuCollectiveDir, IalltoallBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallID] = ialltoallInfo

	alltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallvID,
			BinName: AlltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, AlltoallvBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[AlltoallvID] = alltoallvInfo

	ialltoallvInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallvID,
			BinName: IalltoallvBinName,
			BinPath: filepath.Join(osuCollectiveDir, IalltoallvBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallvID] = ialltoallvInfo

	alltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    AlltoallwID,
			BinName: AlltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, AlltoallwBinName),
		},
		OutputParser: nil,
	}
	m[AlltoallwID] = alltoallwInfo

	ialltoallwInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IalltoallwID,
			BinName: IalltoallwBinName,
			BinPath: filepath.Join(osuCollectiveDir, IalltoallwBinName),
		},
		OutputParser: nil,
	}
	m[IalltoallwID] = ialltoallwInfo

	barrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    BarrierID,
			BinName: BarrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, BarrierBinName),
		},
		OutputParser: nil,
	}
	m[BarrierID] = barrierInfo

	ibarrierInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IbarrierID,
			BinName: IbarrierBinName,
			BinPath: filepath.Join(osuCollectiveDir, IbarrierBinName),
		},
		OutputParser: nil,
	}
	m[IbarrierID] = ibarrierInfo

	bcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    BcastID,
			BinName: BcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, BcastBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[BcastID] = bcastInfo

	ibcastInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IbcastID,
			BinName: IbcastBinName,
			BinPath: filepath.Join(osuCollectiveDir, IbcastBinName),
		},
		OutputParser: nil,
	}
	m[IbcastID] = ibcastInfo

	gatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    GatherID,
			BinName: GatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, GatherBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[GatherID] = gatherInfo

	igatherInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IgatherID,
			BinName: IgatherBinName,
			BinPath: filepath.Join(osuCollectiveDir, IgatherBinName),
		},
		OutputParser: nil,
	}
	m[IgatherID] = igatherInfo

	gathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    GathervID,
			BinName: GathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, GathervBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[GathervID] = gathervInfo

	igathervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IgathervID,
			BinName: IgathervBinName,
			BinPath: filepath.Join(osuCollectiveDir, IgathervBinName),
		},
		OutputParser: nil,
	}
	m[IgathervID] = igathervInfo

	reduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ReduceID,
			BinName: ReduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, ReduceBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ReduceID] = reduceInfo

	ireduceInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IreduceID,
			BinName: IreduceBinName,
			BinPath: filepath.Join(osuCollectiveDir, IreduceBinName),
		},
		OutputParser: nil,
	}
	m[IreduceID] = ireduceInfo

	reduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ReduceScatterID,
			BinName: ReduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, ReduceScatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ReduceScatterID] = reduceScatterInfo

	ireduceScatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IreduceScatterID,
			BinName: IreduceScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, IreduceScatterBinName),
		},
		OutputParser: nil,
	}
	m[IreduceScatterID] = ireduceScatterInfo

	scatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ScatterID,
			BinName: ScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, ScatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ScatterID] = scatterInfo

	iscatterInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IscatterID,
			BinName: IscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, IscatterBinName),
		},
		OutputParser: nil,
	}
	m[IscatterID] = iscatterInfo

	scattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    ScattervID,
			BinName: ScatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, ScatterBinName),
		},
		OutputParser: results.ParseOutputFile,
	}
	m[ScattervID] = scattervInfo

	iscattervInfo := SubBenchmark{
		AppInfo: app.Info{
			Name:    IscattervID,
			BinName: IscatterBinName,
			BinPath: filepath.Join(osuCollectiveDir, IscatterBinName),
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
	b.App.Source.URL = cfg.URL

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
