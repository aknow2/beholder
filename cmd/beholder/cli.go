package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aknow2/beholder/internal/app"
	"github.com/aknow2/beholder/internal/summary"
)

// Version is injected at build time via -X ldflags
var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "record":
		recordCmd(args)
	case "events":
		eventsCmd(args)
	case "summary":
		summaryCmd(args)
	case "reset":
		resetCmd(args)
	case "version", "--version", "-v":
		versionCmd()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func recordCmd(args []string) {
	fs := flag.NewFlagSet("record", flag.ExitOnError)
	configPath := fs.String("config", "~/.beholder/config.yaml", "path to config file")
	oneshoot := fs.Bool("oneshot", false, "record a single event and exit")
	_ = fs.Parse(args)

	ctx := context.Background()
	appInstance, err := app.NewApp(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}
	defer appInstance.Close()

	if *oneshoot {
		event, err := appInstance.RecordOnce(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "record error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("recorded: id=%s category=%s confidence=%.2f status=%s\n", event.ID, event.CategoryName, event.Confidence, event.Status)
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nreceived interrupt signal, stopping...")
		cancel()
	}()

	fmt.Printf("starting scheduler (interval: %d minutes)\n", appInstance.Config.Scheduler.IntervalMinutes)
	fmt.Println("press Ctrl+C to stop")

	if err := appInstance.StartScheduler(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "scheduler error: %v\n", err)
		os.Exit(1)
	}
}

func eventsCmd(args []string) {
	fs := flag.NewFlagSet("events", flag.ExitOnError)
	configPath := fs.String("config", "~/.beholder/config.yaml", "path to config file")
	dateStr := fs.String("date", time.Now().Format("2006-01-02"), "date (YYYY-MM-DD)")
	_ = fs.Parse(args)

	date, err := time.ParseInLocation("2006-01-02", *dateStr, time.Local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date: %v\n", err)
		os.Exit(1)
	}

	appInstance, err := app.NewApp(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}
	defer appInstance.Close()

	events, err := appInstance.ListEventsByDate(date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "list error: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("no events")
		return
	}

	for _, e := range events {
		fmt.Printf("%s | category=%s | confidence=%.2f | status=%s\n", e.CapturedAt.Format(time.RFC3339), e.CategoryName, e.Confidence, e.Status)
	}
}

func summaryCmd(args []string) {
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	configPath := fs.String("config", "~/.beholder/config.yaml", "path to config file")
	dateStr := fs.String("date", time.Now().Format("2006-01-02"), "date (YYYY-MM-DD)")
	format := fs.String("format", "text", "output format: text|markdown")
	_ = fs.Parse(args)

	date, err := time.ParseInLocation("2006-01-02", *dateStr, time.Local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date: %v\n", err)
		os.Exit(1)
	}

	appInstance, err := app.NewApp(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}
	defer appInstance.Close()

	events, err := appInstance.ListEventsByDate(date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "list error: %v\n", err)
		os.Exit(1)
	}

	// T022: Remove categoryMap generation, Generate() uses event.CategoryName directly
	dailySummary := summary.Generate(events)

	switch *format {
	case "markdown":
		fmt.Println(dailySummary.FormatMarkdown())
	case "text":
		fmt.Println(dailySummary.FormatText())
	default:
		fmt.Fprintf(os.Stderr, "unknown format: %s\n", *format)
		os.Exit(1)
	}
}

func resetCmd(args []string) {
	fs := flag.NewFlagSet("reset", flag.ExitOnError)
	configPath := fs.String("config", "~/.beholder/config.yaml", "path to config file")
	dateStr := fs.String("date", time.Now().Format("2006-01-02"), "date (YYYY-MM-DD)")
	_ = fs.Parse(args)

	date, err := time.ParseInLocation("2006-01-02", *dateStr, time.Local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("This will delete events for %s. Continue? [y/N]: ", date.Format("2006-01-02"))
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %v\n", err)
		os.Exit(1)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("cancelled")
		return
	}

	appInstance, err := app.NewApp(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}
	defer appInstance.Close()

	deleted, err := appInstance.DeleteEventsByDate(date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "delete error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("deleted %d events for %s\n", deleted, date.Format("2006-01-02"))
}

func versionCmd() {
	fmt.Printf("beholder version %s\n", Version)
}

func printUsage() {
	fmt.Println("Usage: beholder <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  record   start scheduled recording (use --oneshot for single capture)")
	fmt.Println("  events   list events for a date")
	fmt.Println("  summary  generate daily summary report")
	fmt.Println("  reset    delete events for a date (requires confirmation)")
	fmt.Println("  version  display version")
	fmt.Println("Options:")
	fmt.Println("  --config <path>      config file path (default: ~/.beholder/config.yaml)")
	fmt.Println("  --date <YYYY-MM-DD>  date for events/summary (default: today)")
	fmt.Println("  --format <type>      output format for summary: text|markdown (default: text)")
}
