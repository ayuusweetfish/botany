# Botany: Judges

A judge is an ordinary submission compiled with the same compiler and flags as
participant submissions.

The judge receives the names and desired log outputs through command line args.
It creates child processes, redirect their stderr to log files, interact with
them in any way appropriate and writes the report to stdout. Its own stderr is
used to log the judge's activity and report internal errors. For errors to be
reported with a "System Error" status in the system, the judge returns any non-
zero status.

The judge is agnostic to updates to ratings and performance data, which are done
with contest scripts that read reports produced by the judge.
