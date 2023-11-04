# Command line tool (sqlgen)

## Syntax

Use the following syntax to run `sqlgen` commands from your terminal window:

```bash
sqlgen [command] [flags]
```

where `command`, and `flags` are:

- `command`: Specifies the operation that you want to perform.
- `flags`: Specifies optional flags.

## Cheat Sheet

- Create a configuration file `sqlgen.yml` in current directory.

```bash
sqlgen init
```

- Check `sqlgen` version.

```bash
sqlgen version
```

- Generate necessary code for specific location, glob is support.

```bash
# sqlgen generate [source]
sqlgen generate ./examples/*.go

# or
sqlgen -c config.yml # this will load `config.yml`

# or
sqlgen # this will load `sqlgen.yml`
```
