package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/agent"
)

var (
	REVISION           string
	targetHost         string
	profileFile        string
	hostAdvertise      string
	tlsCertificatePath string
	tlsKeyPath         string
	useTLS             bool
	exitStatusOnFail   bool
	noLoad             bool
	promOut            string
	showVersion        bool
	progress           bool
	parallelism        int
)

func init() {
	certs, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	agent.DefaultTLSConfig.ClientCAs = certs
	agent.DefaultTLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	agent.DefaultTLSConfig.MinVersion = tls.VersionTLS12
	agent.DefaultTLSConfig.InsecureSkipVerify = false

	isAdmin := false
	flag.StringVar(&targetHost, "target", os.Getenv("BENCHMARKER_TARGET_HOST"), "ex: 127.0.0.1:9292")
	flag.StringVar(&profileFile, "profile", "", "ex: cpu.out")
	flag.StringVar(&hostAdvertise, "host-advertise", "local.t.isucon.dev", "hostname to advertise against target")
	flag.StringVar(&tlsCertificatePath, "tls-cert", "../secrets/cert.pem", "path to TLS certificate for a push service")
	flag.StringVar(&tlsKeyPath, "tls-key", "../secrets/key.pem", "path to private key of TLS certificate for a push service")
	flag.BoolVar(&exitStatusOnFail, "exit-status", false, "set exit status non-zero when a benchmark result is failing")
	flag.BoolVar(&noLoad, "no-load", false, "exit on finished prepare")
	flag.StringVar(&promOut, "prom-out", "", "Prometheus textfile output path")
	flag.BoolVar(&showVersion, "version", false, "show version and exit 1")
	flag.IntVar(&parallelism, "parallelism", 20, "parallelism count")
	flag.BoolVar(&progress, "progress", false, "show score in progress")
	flag.BoolVar(&isAdmin, "admin", false, "administrator mode")

	timeoutDuration := ""
	flag.StringVar(&timeoutDuration, "timeout", "10s", "request timeout duration")

	flag.Parse()

	timeout, err := time.ParseDuration(timeoutDuration)
	if err != nil {
		panic(err)
	}
	agent.DefaultRequestTimeout = timeout

	if !isAdmin {
		AdminLogger = log.New(&Blackhole{}, "", log.Lmicroseconds)
	}
}

func checkError(err error) (critical bool, timeout bool, deduction bool) {
	critical = isCritical(err)
	timeout = isTimeout(err)
	deduction = !timeout && isDeduction(err)

	return
}

func sendResult(s *Scenario, result *isucandar.BenchmarkResult, finish bool) bool {
	passed := true
	reason := "pass"
	errors := result.Errors.All()

	result.Score.Set(ScoreSignup, 0)
	result.Score.Set(ScoreLogin, 1)
	result.Score.Set(ScoreCreateSchedule, 10)
	result.Score.Set(ScoreCreateReservation, 1)

	scoreRaw := result.Score.Sum()
	deduction := int64(0)
	timeoutCount := int64(0)

	for tag, count := range result.Score.Breakdown() {
		AdminLogger.Printf("SCORE: %s: %d", tag, count)
	}

	for _, err := range errors {
		isCritical, isTimeout, isDeduction := checkError(err)

		switch true {
		case isCritical:
			passed = false
			reason = "fail: critical"
		case isTimeout:
			timeoutCount++
		case isDeduction:
			if isValidation(err) {
				deduction += 50
			} else {
				deduction++
			}
		}
	}
	deductionTotal := deduction + timeoutCount/10

	score := scoreRaw - deductionTotal
	if score <= 0 && passed {
		passed = false
		reason = "fail: score"
	}

	ContestantLogger.Printf("score: %d(%d - %d) : %s", score, scoreRaw, deductionTotal, reason)
	ContestantLogger.Printf("deduction: %d / timeout: %d", deduction, timeoutCount)

	return passed
}

func writePromFile(promTags []string) {
	if len(promOut) == 0 {
		return
	}

	promOutNew := fmt.Sprintf("%s.new", promOut)
	err := ioutil.WriteFile(promOutNew, []byte(strings.Join(promTags, "")), 0644)
	if err != nil {
		AdminLogger.Printf("Failed to write prom file: %s", err)
		return
	}
	err = os.Rename(promOutNew, promOut)
	if err != nil {
		AdminLogger.Printf("Failed to write prom file: %s", err)
		return
	}

}

func main() {
	AdminLogger.Printf("ISUCON11 Prior benchmarker %s", REVISION)

	if showVersion {
		fmt.Println(REVISION)
		os.Exit(1)
	}

	if profileFile != "" {
		fs, err := os.Create(profileFile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(fs)
		defer pprof.StopCPUProfile()
	}
	if targetHost == "" {
		targetHost = "localhost"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := NewScenario()
	scheme := "http"
	if useTLS {
		scheme = "https"
	}
	s.BaseURL = fmt.Sprintf("%s://%s/", scheme, targetHost)
	s.NoLoad = noLoad
	s.Parallelism = int32(parallelism)

	b, err := isucandar.NewBenchmark(
		isucandar.WithLoadTimeout(70 * time.Second),
	)
	if err != nil {
		panic(err)
	}

	errorCount := int64(0)
	b.OnError(func(err error, step *isucandar.BenchmarkStep) {
		critical, timeout, deduction := checkError(err)
		if timeout {
			return
		}

		if critical || (deduction && atomic.AddInt64(&errorCount, 1) >= 100) {
			step.Cancel()
		}

		ContestantLogger.Printf("ERR: %v", err)
	})

	b.AddScenario(s)

	wg := sync.WaitGroup{}

	b.Load(func(ctx context.Context, step *isucandar.BenchmarkStep) error {
		if s.NoLoad || !progress {
			return nil
		}

		wg.Add(1)
		defer wg.Done()

		for {
			// 途中経過を3秒毎に送信
			timer := time.After(3 * time.Second)
			sendResult(s, step.Result(), false)

			select {
			case <-timer:
			case <-ctx.Done():
				return nil
			}
		}
	})

	result := b.Start(ctx)

	wg.Wait()

	if !sendResult(s, result, true) && exitStatusOnFail {
		os.Exit(1)
	}
}
