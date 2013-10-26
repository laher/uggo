uggo
=====

Ungraceful Gnu-ification for Go

This helps Go commandline apps to behave a bit more 'coreutils'-like, or 'Gnu-ish'.

The main feature being to treat options such as `-lah` are treated as `-l -a -h`, whereas `--lah` is treated as-is, as a single option `--lah`
