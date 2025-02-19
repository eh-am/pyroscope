package adhoc

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/pyroscope-io/pyroscope/pkg/agent/pyspy"
	"github.com/pyroscope-io/pyroscope/pkg/agent/rbspy"
	"github.com/pyroscope-io/pyroscope/pkg/agent/types"
	"github.com/pyroscope-io/pyroscope/pkg/agent/upstream/direct"
	"github.com/pyroscope-io/pyroscope/pkg/config"
	"github.com/pyroscope-io/pyroscope/pkg/exec"
	"github.com/pyroscope-io/pyroscope/pkg/exporter"
	"github.com/pyroscope-io/pyroscope/pkg/storage"
)

func newConnect(cfg *config.Adhoc, st *storage.Storage, logger *logrus.Logger) (runner, error) {
	spyName := cfg.SpyName
	if cfg.Pid == -1 {
		if spyName != "" && spyName != "ebpfspy" {
			return nil, fmt.Errorf("pid -1 can only be used with ebpfspy")
		}
		spyName = "ebpfspy"
	}
	if err := exec.PerformChecks(spyName); err != nil {
		return nil, err
	}

	upstream := direct.New(st, exporter.MetricsExporter{})

	// if the sample rate is zero, use the default value
	sampleRate := uint32(types.DefaultSampleRate)
	if cfg.SampleRate != 0 {
		sampleRate = uint32(cfg.SampleRate)
	}

	// TODO: this is somewhat hacky, we need to find a better way to configure agents
	pyspy.Blocking = cfg.PyspyBlocking
	rbspy.Blocking = cfg.RbspyBlocking

	return &exec.Connect{
		Logger:             logger,
		Upstream:           upstream,
		SpyName:            spyName,
		ApplicationName:    exec.CheckApplicationName(logger, cfg.ApplicationName, spyName, []string{}),
		SampleRate:         sampleRate,
		DetectSubprocesses: cfg.DetectSubprocesses,
		Pid:                cfg.Pid,
	}, nil
}
