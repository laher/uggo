uggo
=====

Ungraceful Gnu-ification for Go

This helps Go commandline apps to behave a bit more 'coreutils'-like, or 'Gnu-ish'.

## Initial features

 * treat options such as `-lah` are treated as `-l -a -h`, whereas `--lah` is treated as-is, as a single option `--lah`
 * detect whether STDIN is being piped from another process

## TODO
 
 * More comprehensive flagset wrapper. e.g. definition of long & short options
