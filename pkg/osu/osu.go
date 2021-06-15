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

func getPt2PtSubBenchmarks(basedir string, wp *workspace.Config, m map[string]app.Info) {
	p2pRelativePath := filepath.Join("libexec", "osu-micro-benchmarks", "mpi", "pt2pt")
	osuPt2Pt2Dir := filepath.Join(wp.InstallDir, basedir, p2pRelativePath)

	bwInfo := app.Info{
		Name:    bwBinID,
		BinName: bwBinName,
		BinPath: filepath.Join(osuPt2Pt2Dir, bwBinName),
	}
	m[bwBinID] = bwInfo

	latencyInfo := app.Info{
		Name:    latencyBinID,
		BinName: LatencyBinName,
		BinPath: filepath.Join(osuPt2Pt2Dir, LatencyBinName),
	}
	m[latencyBinID] = latencyInfo
}

func getCollectiveSubBenchmarks(basedir string, wp *workspace.Config, m map[string]app.Info) {
	collectiveRelativePath := filepath.Join("libexec", "osu-micro-benchmarks", "mpi", "collective")
	osuCollectiveDir := filepath.Join(wp.InstallDir, basedir, collectiveRelativePath)

	allgatherInfo := app.Info{
		Name:    allgatherBinID,
		BinName: allgatherBinName,
		BinPath: filepath.Join(osuCollectiveDir, allgatherBinName),
	}
	m[allgatherBinID] = allgatherInfo

	iallgatherInfo := app.Info{
		Name:    iallgatherBinID,
		BinName: iallgatherBinName,
		BinPath: filepath.Join(osuCollectiveDir, iallgatherBinName),
	}
	m[iallgatherBinID] = iallgatherInfo

	allgathervInfo := app.Info{
		Name:    allgathervBinID,
		BinName: allgathervBinName,
		BinPath: filepath.Join(osuCollectiveDir, allgathervBinName),
	}
	m[allgathervBinID] = allgathervInfo

	iallgathervInfo := app.Info{
		Name:    iallgathervBinID,
		BinName: iallgathervBinName,
		BinPath: filepath.Join(osuCollectiveDir, iallgathervBinName),
	}
	m[iallgathervBinID] = iallgathervInfo

	allreduceInfo := app.Info{
		Name:    allreduceBinID,
		BinName: allreduceBinName,
		BinPath: filepath.Join(osuCollectiveDir, allreduceBinName),
	}
	m[allreduceBinID] = allreduceInfo

	iallreduceInfo := app.Info{
		Name:    iallreduceBinID,
		BinName: iallreduceBinName,
		BinPath: filepath.Join(osuCollectiveDir, iallreduceBinName),
	}
	m[iallreduceBinID] = iallreduceInfo

	alltoallInfo := app.Info{
		Name:    alltoallBinID,
		BinName: alltoallBinName,
		BinPath: filepath.Join(osuCollectiveDir, alltoallBinName),
	}
	m[alltoallBinID] = alltoallInfo

	ialltoallInfo := app.Info{
		Name:    ialltoallBinID,
		BinName: ialltoallBinName,
		BinPath: filepath.Join(osuCollectiveDir, ialltoallBinName),
	}
	m[ialltoallBinID] = ialltoallInfo

	alltoallvInfo := app.Info{
		Name:    alltoallvBinID,
		BinName: alltoallvBinName,
		BinPath: filepath.Join(osuCollectiveDir, alltoallvBinName),
	}
	m[alltoallvBinID] = alltoallvInfo

	ialltoallvInfo := app.Info{
		Name:    ialltoallvBinID,
		BinName: ialltoallvBinName,
		BinPath: filepath.Join(osuCollectiveDir, ialltoallvBinName),
	}
	m[ialltoallvBinID] = ialltoallvInfo

	alltoallwInfo := app.Info{
		Name:    alltoallwBinID,
		BinName: alltoallwBinName,
		BinPath: filepath.Join(osuCollectiveDir, alltoallwBinName),
	}
	m[alltoallwBinID] = alltoallwInfo

	ialltoallwInfo := app.Info{
		Name:    ialltoallwBinID,
		BinName: ialltoallwBinName,
		BinPath: filepath.Join(osuCollectiveDir, ialltoallwBinName),
	}
	m[ialltoallwBinID] = ialltoallwInfo

	barrierInfo := app.Info{
		Name:    barrierBinID,
		BinName: barrierBinName,
		BinPath: filepath.Join(osuCollectiveDir, barrierBinName),
	}
	m[barrierBinID] = barrierInfo

	ibarrierInfo := app.Info{
		Name:    ibarrierBinID,
		BinName: ibarrierBinName,
		BinPath: filepath.Join(osuCollectiveDir, ibarrierBinName),
	}
	m[ibarrierBinID] = ibarrierInfo

	bcastInfo := app.Info{
		Name:    bcastBinID,
		BinName: bcastBinName,
		BinPath: filepath.Join(osuCollectiveDir, bcastBinName),
	}
	m[bcastBinID] = bcastInfo

	ibcastInfo := app.Info{
		Name:    ibcastBinID,
		BinName: ibcastBinName,
		BinPath: filepath.Join(osuCollectiveDir, ibcastBinName),
	}
	m[ibcastBinID] = ibcastInfo

	gatherInfo := app.Info{
		Name:    gatherBinID,
		BinName: gatherBinName,
		BinPath: filepath.Join(osuCollectiveDir, gatherBinName),
	}
	m[gatherBinID] = gatherInfo

	igatherInfo := app.Info{
		Name:    igatherBinID,
		BinName: igatherBinName,
		BinPath: filepath.Join(osuCollectiveDir, igatherBinName),
	}
	m[igatherBinID] = igatherInfo

	gathervInfo := app.Info{
		Name:    gathervBinID,
		BinName: gathervBinName,
		BinPath: filepath.Join(osuCollectiveDir, gathervBinName),
	}
	m[gathervBinID] = gathervInfo

	igathervInfo := app.Info{
		Name:    igathervBinID,
		BinName: igathervBinName,
		BinPath: filepath.Join(osuCollectiveDir, igathervBinName),
	}
	m[igathervBinID] = igathervInfo

	reduceInfo := app.Info{
		Name:    reduceBinID,
		BinName: reduceBinName,
		BinPath: filepath.Join(osuCollectiveDir, reduceBinName),
	}
	m[reduceBinID] = reduceInfo

	ireduceInfo := app.Info{
		Name:    ireduceBinID,
		BinName: ireduceBinName,
		BinPath: filepath.Join(osuCollectiveDir, ireduceBinName),
	}
	m[ireduceBinID] = ireduceInfo

	reduceScatterInfo := app.Info{
		Name:    reduceScatterBinID,
		BinName: reduceScatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, reduceScatterBinName),
	}
	m[reduceScatterBinID] = reduceScatterInfo

	ireduceScatterInfo := app.Info{
		Name:    ireduceScatterBinID,
		BinName: ireduceScatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, ireduceScatterBinName),
	}
	m[ireduceScatterBinID] = ireduceScatterInfo

	scatterInfo := app.Info{
		Name:    scatterBinID,
		BinName: scatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
	}
	m[scatterBinID] = scatterInfo

	iscatterInfo := app.Info{
		Name:    iscatterBinID,
		BinName: iscatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
	}
	m[iscatterBinID] = iscatterInfo

	scattervInfo := app.Info{
		Name:    scattervBinID,
		BinName: scatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, scatterBinName),
	}
	m[scattervBinID] = scattervInfo

	iscattervInfo := app.Info{
		Name:    iscattervBinID,
		BinName: iscatterBinName,
		BinPath: filepath.Join(osuCollectiveDir, iscatterBinName),
	}
	m[iscattervBinID] = iscattervInfo
}

func getSubBenchmarks(cfg *benchmark.Config, wp *workspace.Config, basedir string) map[string]app.Info {
	m := make(map[string]app.Info)
	getCollectiveSubBenchmarks(basedir, wp, m)
	getPt2PtSubBenchmarks(basedir, wp, m)
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
	err = mpiInfo.Load()
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
			if fsutil.FileExists(benchmarkInfo.BinPath) {
				log.Printf("\t%s exists", benchmarkInfo.BinPath)
				installInfo.SubBenchmarks = append(installInfo.SubBenchmarks, benchmarkInfo)
			} else {
				log.Printf("\t%s does not exist", benchmarkInfo.BinPath)
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
