package config

import (
	"os"
	"time"

	"github.com/criyle/go-judge/pkg/envexec"
	"github.com/koding/multiconfig"
)

// Config defines executor server configuration
type Config struct {
	// container
	ContainerInitPath  string `flagUsage:"container init path"`
	PreFork            int    `flagUsage:"control # of the prefork workers"`
	TmpFsParam         string `flagUsage:"tmpfs mount data (only for default mount with no mount.yaml)" default:"size=16m,nr_inodes=4k"`
	NetShare           bool   `flagUsage:"share net namespace with host"`
	MountConf          string `flagUsage:"specifies mount configuration file" default:"mount.yaml"`
	SeccompConf        string `flagUsage:"specifies seccomp filter" default:"seccomp.yaml"`
	Parallelism        int    `flagUsage:"control the # of concurrency execution" default:"4"`
	CgroupPrefix       string `flagUsage:"control cgroup prefix" default:"executor_server"`
	ContainerCredStart int    `flagUsage:"control the start uid&gid for container" default:"10000"`

	// file store
	SrcPrefix string `flagUsage:"specifies directory prefix for source type copyin"`
	Dir       string `flagUsage:"specifies directory to store file upload / download (in memory by default)"`

	// runner limit
	TimeLimitCheckerInterval time.Duration `flagUsage:"specifies time limit checker interval" default:"100ms"`
	ExtraMemoryLimit         *envexec.Size `flagUsage:"specifies extra memory buffer for check memory limit" default:"16k"`
	OutputLimit              *envexec.Size `flagUsage:"specifies POSIX rlimit for output for each command" default:"256m"`
	Cpuset                   string        `flagUsage:"control the usage of cpuset for all containerd process"`
	EnableCPURate            bool          `flagUsage:"enable cpu cgroup rate control"`
	CPUCfsPeriod             time.Duration `flagUsage:"set cpu.cfs_period" default:"100ms"`

	// server config
	HTTPAddr      string `flagUsage:"specifies the http binding address" default:":5050"`
	EnableGRPC    bool   `flagUsage:"enable gRPC endpoint"`
	GRPCAddr      string `flagUsage:"specifies the grpc binding address" default:":5051"`
	AuthToken     string `flagUsage:"bearer token auth for REST / gRPC"`
	EnableDebug   bool   `flagUsage:"enable debug endpoint"`
	EnableMetrics bool   `flagUsage:"enable promethus metrics endpoint"`

	// logger config
	Release bool `flagUsage:"release level of logs"`
	Silent  bool `flagUsage:"do not print logs"`
}

// Load loads config from flag & environment variables
func (c *Config) Load() error {
	cl := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{
			Prefix:    "ES",
			CamelCase: true,
		},
		&multiconfig.FlagLoader{
			CamelCase: true,
			EnvPrefix: "ES",
		},
	)
	if os.Getpid() == 1 {
		c.Release = true
	}
	return cl.Load(c)
}
