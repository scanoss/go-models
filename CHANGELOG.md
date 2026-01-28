# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-01-29

### Changed
- `GetComponent` now bypasses version resolution when the requirement is a fixed version (no range operators)
  - Returns directly without database lookup for exact versions
  - Covers operators from npm (`^~<>=*|`), pip (`~=!<>=`), maven/nuget (`[](),`), cargo, composer, and gems

[Unreleased]: https://github.com/scanoss/go-models/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/scanoss/go-models/releases/tag/v0.3.0