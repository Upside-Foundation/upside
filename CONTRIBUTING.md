# Contributing

Thank you for your interest in contributing to Upside. This repository holds `upsided`, the node software for the
Upside network: an EVM-compatible Layer 1 built on the Cosmos SDK that serves as the settlement layer for the Upside
prediction network. For chat and quick questions, join the [Upside Discord server](https://discord.gg/wc7C5AZhh3). For other questions, open a GitHub Discussion or issue.

## Code of Conduct

All contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Requests involving no code changes (e.g. conventions, processes or network parameters)

Requests involving no code changes (e.g. conventions, processes or network parameters) should be posted as GitHub Discussions or issues to allow for early feedback before e.g. being submitted for formal network governance decisions.

## Feature requests and bug reports

Feature requests and bug reports should be posted as GitHub issues.
In an issue, please describe what you did, what you expected, and what happened instead.

If you think that you have identified an issue with Upside that might compromise
its users' security, please do not open a public issue on GitHub. Instead,
we ask you to refer to [security policy](SECURITY.md).

## Consensus-breaking and large structural code changes

Describe the proposed behavior, compatibility requirements, and validation plan in an issue before implementing consensus or substantial structural changes. Consensus changes also need an upgrade and migration plan.

## Working on issues
There are several ways to identify an area where you can contribute to Upside:

* You can open a GitHub issue, either with a specific contribution in mind or in general by saying "I want to help!".
* Occasionally, some issues on GitHub may be labelled with `help wanted` or `good first issue` tags.

We use the variation of the "fork and pull" model where contributors push changes to their personal fork and create pull requests to bring those changes into the source repository.
The code and comments should follow [Effective Go Guide](https://golang.org/doc/effective_go.html) (and [Uber Style Guide](https://github.com/uber-go/guide/blob/master/style.md)). Many of the code style rules are captured by `go fmt`, `go vet`, `golangci-lint` (`make lint`), and other tools, so we recommend [setting up your editor to do formatting and lint-checking for you](https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins).

Once you identified an issue to work on, this is the summary of your basic steps:

* Fork Upside's repository under your GitHub account.

* Clone your fork locally on your machine.

* Post a comment in the issue to say that you are working on it, so that other people do not work on the same issue.

* Create a local branch on your machine by `git checkout -b branch_name`.

* Commit your changes to your own fork.

* Include tests that cover all non-trivial code.

* Check you are working on the latest version on main in the official repository. If not, please pull the official repository's main (upstream) into your fork's main branch, and rebase your committed changes or replay your stashed changes in your branch over the latest changes in the upstream version.

* Run all tests locally and make sure they pass.

* If your changes are of interest to other developers, please make corresponding changes in the official documentation and the changelog.

* Push your changes to your fork's branch and open the pull request to the repository main branch.

* In the pull request, add a clear description of the problem your changes solve, and add the following statement to confirm that your contribution is your own original work: "I hereby certify that my contribution is in accordance with the Developer Certificate of Origin (https://developercertificate.org/)."

* The reviewer will either accept and merge your pull request, or leave comments requesting changes via the GitHub PR interface (you should then make changes by pushing directly to your existing PR branch).

* Pull requests need review and approval from a maintainer before they merge. Changes that break consensus also need a governance software upgrade; see [Releases](#releases).

### Changelog

Every non-trivial PR must update the [CHANGELOG.md](./CHANGELOG.md).

The Changelog is *not* a record of what Pull Requests were merged;
the commit history already shows that. The Changelog is a notice to the user
about how their expectations of the software should be modified. 
It is part of the UX of a release and is a *critical* user facing integration point.
The Changelog must be clean, inviting, and readable, with concise, meaningful entries. 
Entries must be semantically meaningful to users. If a change takes multiple
Pull Requests to complete, it should likely have only a single entry in the
Changelog describing the net effect to the user.

When writing Changelog entries, ensure they are targeting users of the software,
not fellow developers. Developers have much more context and care about more
things than users do. Changelogs are for users. 

Changes are grouped under a heading for each release, with `## Unreleased` for changes that
are not yet released. Within a release, entries are split between these sections, as needed:
Features, State machine breaking, Build, and Documentation.

Start each entry with the part of the code that changed, such as `x/upside`, followed by a colon
and a description of the change:
```
- `x/upside`: Some description about the change.
```

## Releases

Our release process is as follows:

1. The [changelog](#changelog) is updated to summarize all changes in the release.
2. Changes that break consensus ship with an upgrade handler registered under the upgrade's name in
   [`app/upgrades.go`](app/upgrades.go). They reach mainnet only through a governance software upgrade
   proposal that names the upgrade, the upgrade height, the source commit, and the binary checksums.
3. The release commit is tagged `vX.Y.Z`.
4. The mainnet release binary is built reproducibly with
   [`scripts/build-release.sh`](scripts/build-release.sh) (pinned Go and image; see
   [Build and Install](docs/getting-started/install.md#reproduce-a-release-with-docker)), and its
   SHA-256 checksum is published with the release and in the upgrade proposal. The Nix build
   ([`scripts/release.sh`](scripts/release.sh)) uses a different Go version and does not produce the same
   bytes.
5. The GitHub release notes are written from the changelog, with links to any upgrade instructions.

### Developer Certificate of Origin
All contributions to this project are subject to the terms of the Developer Certificate of Origin, available [here](https://developercertificate.org/) and reproduced below:

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
