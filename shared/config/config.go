package config

import (
	"fmt"
	"os/exec"
	"strconv"
)

const (
	maxSize = 128 * 1024
)

type config struct {
	// Input items are terminated by a null character instead of by whitespace, and the quotes and backslash are not special (every character is taken literally).
	// Disables the end of file string, which is treated like any other argument. Useful when input items might contain white space, quote marks, or backslashes.
	// The GNU find -print0 option produces input suitable for this mode.
	NullSeparator bool

	// Read items from file instead of standard input. If you use this option, stdin remains unchanged when commands are run. Otherwise, stdin is redirected
	// from /dev/null.
	File string

	// Input items are terminated by the specified character. The specified delimiter may be a single character, a C-style character escape such as \n, or an oc‐
	// tal or hexadecimal escape code. Octal and hexadecimal escape codes are understood as for the printf command.  Multibyte characters are not supported.
	// When processing the input, quotes and backslash are not special; every character in the input is taken literally. The -d option disables any end-of-file
	// string, which is treated like any other argument. You can use this option when the input consists of simply newline-separated items, although it is almost
	// always better to design your program to use --null where this is possible.
	Delimiter string

	// Set the end of file string to eof-str.  If the end of file string occurs as a line of input, the rest of the input is ignored
	EOL bool

	// Use at most max-args arguments per command line. Fewer than max-args arguments will be used if the size (see the -s option) is exceeded, unless the -x op‐
	// tion is given, in which case xargs will exit.
	MaxArgs int

	// Use at most max-chars characters per command line, including the command and initial-arguments and the terminating nulls at the ends of the argument
	// strings. The largest allowed value is system-dependent, and is calculated as the argument length limit for exec, less the size of your environment, less
	// 2048 bytes of headroom. If this value is more than 128KiB, 128KiB is used as the default value; otherwise, the default value is the maximum. 1KiB is 1024
	// bytes.
	Size int

	// Exit if the size (see the -s option) is exceeded.
	Exit bool

	// Print the command line on the standard error output before executing it.
	Verbose bool
}

var configuration config

func Config() *config {
	return &configuration
}

func Init() (err error) {
	if configuration.Size, err = size(); err != nil {
		return
	}
	return
}

func size() (int, error) {
	cmd := exec.Command("getconf", "ARG_MAX")
	argMax, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("getconf %w", err)
	}

	if argMax[len(argMax)-1] == '\n' {
		argMax = argMax[:len(argMax)-1]
	}

	max, err := strconv.Atoi(string(argMax))
	if err != nil {
		return 0, err
	}

	cmd = exec.Command("sh", "-c", "printf \"%s\" \"$(env)\" | wc -c")
	env, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("sh -c printf \"%%s\" \"$(env)\" | wc -c %w", err)
	}

	if env[len(env)-1] == '\n' {
		env = env[:len(env)-1]
	}

	envsize, err := strconv.Atoi(string(env))
	if err != nil {
		return 0, err
	}

	result := (max - envsize - 2048)
	if result > maxSize {
		result = maxSize
	}

	return result, nil
}
