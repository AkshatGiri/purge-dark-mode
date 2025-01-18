# Purge Dark Mode

This repo contains a cli script written in go to remove all the `dark:` tailwind classes from all files within a directory.

## Buy Why???

Because I started a laravel project and picked support for dark mode only to realize I don't want to support two different themes. So I decided to remove all the `dark:` classes from the project.

It's not only laravel, pretty much every framework provides templates with theme support and it's incredibly hard to resist the templtation of a dark and light theme when setting up a new project.

So this script is to remove all those pesky little `dark:` theme tailwind classes from your project.

## How to use

1. Clone the repo
2. Run `go build -o purge-dark-mode main.go`
3. Run `./purge-dark-mode -dir=/path/to/your/directory` ( The defaul is set to the currnet working directory, so be careful running it in this directory since there is a `dark:` in the go code that will be removed. )

Cli Flags

- `-dir` : The directory to search for files. Default is the current working directory.
- `-dry-run` : by setting this flag no actual changes will be made to the files. You can view the logs to make sure everything looks good before running the script without this flag.
- `log-level`: Set the log level. Options are `debug`, `info`, `warn`, `error`. Default is `info`.

By default the gitignored files are ignored by reading the top level `.gitignore` file. Currently multiple .gitignore files are not supported. As of now there is no way to turn this off. The .git directory is also skipped by default.

> [!TIP]
> You should git commit your changes before running the script. This way you can easily review the chagnes and revert them if something goes wrong.

## TODO

- Support multiple `.gitignore` files in children directory.
- Ignore .css files. It's possible to have `dark:` string in css files, we should probably not touch those. So having a way in the cli to provide certain file extensions to ignore would be nice. As a workaround for now, before running the script those file extensions can be added to the .gitignore file and removed after running the script.
