# Contributing to Proji

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:

-   Reporting a bug
-   Discussing the current state of the code
-   Submitting a fix
-   Proposing new features
-   Becoming a maintainer

By participating to this project, you agree to abide our [code of
conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`proji` is written in [Go](https://golang.org/).

Prerequisites:

- `make`
- [Go 1.14+](https://golang.org/doc/install)

Clone `proji` anywhere:

```sh
$ git clone https://github.com/nikoksr/proji
```

Install the build and lint dependencies:

```sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

## Use a consistent coding style

Before commiting any changes, please run:

```sh
$ make fmt
```

Which runs `gofmt` and `goimports` and enforces a consistent coding style.

## Test your change

You can create a branch for your changes and try to build from the source as you go:

```sh
$ make build
```

When you are satisfied with the changes, we suggest you run:

```sh
$ make ci
```

Which runs all the linters and tests.


## We Develop with Github

We use github to host code, to track issues and feature requests, as well as accept pull requests.

## We Use [Github Flow](https://guides.github.com/introduction/flow/index.html), So All Code Changes Happen Through Pull Requests

Pull requests are the best way to propose changes to the codebase (we use [Github Flow](https://guides.github.com/introduction/flow/index.html)). We actively welcome your pull requests:

1.  Fork the repo and create your branch from `master`.
2.  If you've added code that should be tested, add tests.
3.  If you've added/removed/updated cli commands, update the documentation.
4.  Ensure the test suite passes.
5.  Make sure your code lints.
6.  Issue that pull request!

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issues](https://github.com/nikoksr/proji/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/nikoksr/proji/issues/new/choose); it's that easy!

## Write bug reports with detail, background, and sample code

Use the [bug report template](https://github.com/nikoksr/proji/blob/master/.github/ISSUE_TEMPLATE/bug_report.md) to report a bug using a Github's issues.

## Request features or changes

Use the [feature request template](https://github.com/nikoksr/proji/blob/master/.github/ISSUE_TEMPLATE/feature_request.md) to request a new feature using a Github's issues.

## License

By contributing, you agree that your contributions will be licensed under its MIT License.

## References

This document was adapted from the open-source contribution guidelines for [Facebook's Draft](https://github.com/facebook/draft-js/blob/a9316a723f9e918afde44dea68b5f9f39b7d9b00/CONTRIBUTING.md)
