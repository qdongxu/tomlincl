# tomlincl
A command line tool to parse a toml file and include other toml files.

The **tomlincl** command parse a toml file, insert the contents from other toml
files indicated by the `#!include` directives. **tomlincl** does not parse the
toml semantics. Instead, it merely insert the text lines into the parent toml file.

The `#!include` directives are parsed recursively in the discovered files.

# Usage

1. Add `#!include` directives in the toml file:

> #!include foo*.toml
> #!include *bar.toml
> #!include subdirfoo/*.toml
> #!include subdirbar/*.toml

2. run the **tomlincl** command, the merged lines will print to the stdout:

> % tomlincl \<root toml file\>

