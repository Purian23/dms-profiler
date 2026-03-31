package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"dms-profiler/internal/chromium"
	"dms-profiler/internal/config"
	"dms-profiler/internal/launch"
	"dms-profiler/internal/match"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfgPath := flag.String("config", config.DefaultPath, "path to config.toml")
	dryRun := flag.Bool("dry-run", false, "print argv and exit without launching")
	printCmd := flag.Bool("print-cmd", false, "same as -dry-run")
	flag.Parse()

	if *printCmd {
		*dryRun = true
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <url>\n", os.Args[0])
		return 2
	}
	rawURL := strings.TrimSpace(args[0])
	if rawURL == "" {
		fmt.Fprintln(os.Stderr, "empty URL")
		return 2
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	flat := config.FlattenRules(cfg.Rules)
	rules := make([]match.Rule, 0, len(flat))
	for _, r := range flat {
		rules = append(rules, match.Rule{Match: r.Match, Profile: r.Profile})
	}
	profileSpec := cfg.Default.Profile
	if p, ok := match.FirstPrefixRule(rawURL, rules); ok {
		profileSpec = p
	}
	if strings.TrimSpace(profileSpec) == "" {
		fmt.Fprintln(os.Stderr, "no profile resolved (check [default] and rules)")
		return 1
	}

	profileDir, err := chromium.ResolveProfileDir(cfg.Browser.UserDataDir, profileSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "profile: %v\n", err)
		return 1
	}

	bin, err := launch.LookPath(cfg.Browser.Command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "look up browser: %v\n", err)
		return 1
	}
	argv := launch.Args(bin, profileDir, rawURL)
	if *dryRun {
		fmt.Println(formatArgv(argv))
		return 0
	}

	err = launch.Exec(bin, argv, launch.Environ())
	fmt.Fprintf(os.Stderr, "exec: %v\n", err)
	return 1
}

func formatArgv(argv []string) string {
	b := strings.Builder{}
	for i, a := range argv {
		if i > 0 {
			b.WriteByte(' ')
		}
		if strings.ContainsAny(a, " \t\n\"'\\") {
			b.WriteString(shellQuote(a))
		} else {
			b.WriteString(a)
		}
	}
	return b.String()
}

func shellQuote(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}
