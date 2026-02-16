package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aknow2/beholder/internal/app"
	"github.com/aknow2/beholder/internal/config"
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
	case "init":
		initCmd(args)
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

func initCmd(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	configPath := fs.String("config", "~/.beholder/config.yaml", "path to config file")
	_ = fs.Parse(args)

	resolvedPath, err := config.ResolvePath(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve config path error: %v\n", err)
		os.Exit(1)
	}

	if info, err := os.Stat(resolvedPath); err == nil {
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "config path is a directory: %s\n", resolvedPath)
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		ok, err := promptYesNo(reader, fmt.Sprintf("config already exists at %s. Overwrite? [y/N]: ", resolvedPath), false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			os.Exit(1)
		}
		if !ok {
			fmt.Println("cancelled")
			return
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "config stat error: %v\n", err)
		os.Exit(1)
	}

	defaultCfg, err := config.Default()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load default config error: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	maxWidth, err := promptInt(reader, fmt.Sprintf("image.max_width [%d]: ", defaultCfg.Image.MaxWidth), defaultCfg.Image.MaxWidth, 100, 4096)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %v\n", err)
		os.Exit(1)
	}

	saveImages, err := promptYesNo(reader, fmt.Sprintf("image.save_images [%t] (y/n): ", defaultCfg.Image.SaveImages), defaultCfg.Image.SaveImages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "input error: %v\n", err)
		os.Exit(1)
	}

	defaultCfg.Image.MaxWidth = maxWidth
	defaultCfg.Image.SaveImages = saveImages

	if err := config.Validate(defaultCfg); err != nil {
		fmt.Fprintf(os.Stderr, "config validation error: %v\n", err)
		os.Exit(1)
	}

	if err := config.Write(*configPath, defaultCfg); err != nil {
		fmt.Fprintf(os.Stderr, "write config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("config written to %s\n", resolvedPath)
}

func promptInt(reader *bufio.Reader, prompt string, defaultValue, minValue, maxValue int) (int, error) {
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return defaultValue, nil
		}

		value, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("invalid number, try again")
			continue
		}
		if value < minValue || value > maxValue {
			fmt.Printf("value must be between %d and %d\n", minValue, maxValue)
			continue
		}

		return value, nil
	}
}

func promptYesNo(reader *bufio.Reader, prompt string, defaultValue bool) (bool, error) {
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			return defaultValue, nil
		}
		if line == "y" || line == "yes" {
			return true, nil
		}
		if line == "n" || line == "no" {
			return false, nil
		}
		fmt.Println("please enter y or n")
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
	fmt.Println("  init     create config interactively")
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
