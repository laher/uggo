ugfug
=====

Ungraceful Go Flag Utility for Gnu-ification

This just allows the Go `flag` package to behave a bit more Gnu-ish - a bit more 'coreutils'-like.

The main feature being that options such as `-lah` are treated as `-l -a -h`, whereas `--lah` is treated as-is, as a single option `--lah`
