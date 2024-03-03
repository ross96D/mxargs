package flags

type CmdFlags struct {
	Delimiter string

	NullSeparator bool

	FilePath string

	// Set the end of file string to eof-str.  If the end of file string occurs as a line of input, the rest of the input is ignored
	EOF string

	// Use at most max-args arguments per command line. Fewer than max-args arguments will be used if the size (see the -s option) is exceeded, unless the -x op‐
	// tion is given, in which case xargs will exit.
	MaxArgs int

	// Use at most max-chars characters per command line, including the command and initial-arguments and the terminating nulls at the ends of the argument
	// strings. The largest allowed value is system-dependent, and is calculated as the argument length limit for exec, less the size of your environment, less
	// 2048 bytes of headroom. If this value is more than 128KiB, 128KiB is used as the default value; otherwise, the default value is the maximum. 1KiB is 1024
	// bytes.
	MaxChars uint64

	// Exit if the size (see the -s option) is exceeded.
	Exit bool

	// Print the command line on the standard error output before executing it.
	Verbose bool
}

var Flags = &CmdFlags{}
