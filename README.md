# Heimdallr

[![Build Status][ci-img]][ci]
[![Coverage Status][coverage-img]][coverage]
[![Go Report Card][report-card-img]][report-card]
[![GitHub release][release-img]][release]
[![License][license-img]][license]

[Heimdallr] manages external checks for endpoints in a Kuberneres cluster. It is based
on [Heptio Cruise] but uses Custom Resource Definitions to define Pingdom checks instead
of inferring them from Ingress objects.

[Heimdallr]: https://en.wikipedia.org/wiki/Heimdallr
[Heptio Cruise]: https://github.com/heptiolabs/cruise

[ci-img]: https://travis-ci.org/jeromefroe/heimdallr.svg?branch=master
[ci]: https://travis-ci.org/jeromefroe/heimdallr
[coverage-img]: https://codecov.io/gh/jeromefroe/heimdallr/branch/master/graph/badge.svg
[coverage]: https://codecov.io/gh/jeromefroe/heimdallr
[report-card-img]: https://goreportcard.com/badge/github.com/jeromefroe/heimdallr
[report-card]: https://goreportcard.com/report/github.com/jeromefroe/heimdallr
[release-img]: https://img.shields.io/github/release/jeromefroe/heimdallr.svg
[release]: https://github.com/jeromefroe/heimdallr/releases
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license]: https://raw.githubusercontent.com/jeromefroe/heimdallr/master/LICENSE
