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
