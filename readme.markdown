# ut

`ut` is a command line tool to handle unix timestamps

## Usage

The unix timestamp (`ut`) tool has enough help to get you started:

    $ ut --help

Your local timezone is used, unless manipulated by the flags `--utc` or `--offset`. `--utc` is the equivalent of
`--offset=UTC`.

Other than the help, it has two subcommands to handle timestamps

### Generate

Generates a timestamp for the current time.

    $ ut generate
    1680717044

For more information, run:

    $ ut generate help

### Parse

Parse unix timestamps. You can use the generated value from the `generate` command.

    $ ut --utc parse 1680717044
    2023-04-05 17:50:44 +0000 UTC

For more information, run:

    $ ut parse help

## Inspiration

This tool was inspired by a tool with same name built with Rust, by 
Yoshihito Arih: https://github.com/yoshihitoh/ut-cli

## Author

* [Sergio Moura](https://sergio.moura.com)