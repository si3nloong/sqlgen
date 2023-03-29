## Contribution Guidelines

First off, thanks for taking the time to contribute! ❤️

### Quick Reference

I want to...

_report a bug! ➡️ [open an Issue](#issues)_

_fix a bug! ➡️ [open a Pull Request](#prs)_

_implement a new feature! ➡️ [open an Issue to discuss it first, then a Pull Request](#issues)_

<a name="#src"></a>

### Working on the source code

`cmd/root.go` contains a simple `gopherize` CLI that connects to most of the other source files.

<a name="addex"></a>

### Adding an exercise

The first step is to add the exercise! Name the file `exercise/yourTopic/main.go`, make sure to put in some helpful links, and link to sections of the book in `exercises/yourTopic/README.md`. If you have multiple exercises, you may name the file `exercise/yourTopic/yourTopicN/main.go`. (For example: `exercise/go-routine/chapter1/main.go`)

Next make sure it runs with `go`. The exercise metadata is stored in `exercise.yaml`.

Add the metadata for your exercise in the correct order in the `exercises` array. If you are unsure of the correct ordering, add it at the bottom and ask in your pull request. The exercise metadata should contain the following:

```diff
  ...
+ - name: helloworld
+   type: compile
+   file: exercise/hello_world/main.go
  ...
```

That's all! Feel free to put up a pull request.

<a name="issues"></a>

### Issues

You can open an issue [here](https://github.com/jeffreylean/gopherize/issues/new).
If you're reporting a bug, please include the output of the following commands:

- `go version`
- `ls -la`
- Your OS name and version

<a name="prs"></a>

### Pull Requests

Opening a pull request is as easy as forking the repository and committing your
changes. There's a couple of things to watch out for:

#### Write correct commit messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0-beta.4/)
specification.
This means that you have to format your commit messages in a specific way. Say
you're working on adding a new exercise called `foobar1.go`. You could write
the following commit message:

```
feat: add foobar1.go exercise
```

If you're just fixing a bug, please use the `fix` type:

```
fix(verify): make sure verify doesn't self-destruct
```

The scope within the brackets is optional, but should be any of these:

- `installation` (for the installation script)
- `cli` (for general CLI changes)
- `verify` (for the verification source file)
- `watch` (for the watch functionality source)
- `run` (for the run functionality source)
- `EXERCISENAME` (if you're changing a specific exercise, or set of exercises,
  substitute them here)

When the commit also happens to close an existing issue, link it in the message
body:

```
fix: update foobar

closes #101029908
```

If you're doing simple changes, like updating a book link, use `chore`:

```
chore: update exercise1.go book link
```

If you're updating documentation, use `docs`:

```
docs: add more information to Readme
```

If, and only if, you're absolutely sure you want to make a breaking change
(please discuss this beforehand!), add an exclamation mark to the type and
explain the breaking change in the message body:

```
fix!: completely change verification

BREAKING CHANGE: This has to be done because lorem ipsum dolor
```

### Pull Request Workflow

Once you open a Pull Request, it may be reviewed or labeled (or both) until
the maintainers accept your change. Please be patient, it may take some time
for this to happen!

### License

By contributing to sqlgen, you agree that your contributions will be licensed under its MIT license.
