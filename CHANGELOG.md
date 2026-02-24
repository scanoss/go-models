# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.1] - 2026-02-24
### Changed
- Removed operator check from `GetComponent`
- Optimized `pickOneUrl` in `ComponentService`: replaced `map[*semver.Version]` + `sort` with single-pass max tracking using `GreaterThan`, reducing complexity from O(n log n) to O(n)
- Fixed broken deduplication in `pickOneUrl` caused by using pointer keys in map
- Pre-parsed `v0.0.0` fallback version once instead of on every parse failure

## [0.5.1] - 2026-02-24
### Changed
- Removed operator check from `GetComponent`

## [0.5.0] - 2026-02-24
### Added
- `CheckPurlByNameType` method in `ProjectModel` to check the projects table for entries matching a PURL name and type
- `CheckPurl` method in `ComponentService` to validate and check a PURL string existence

## [0.4.0] - 2026-01-30
### Added
- Added `DBVersionModel` to query the `db_version` table for schema version, package name, release, and creation date
- Added `ErrTableNotFound` sentinel error and `tableExists` helper for backward compatibility with databases that predate certain tables

## [0.3.0] - 2026-01-29
### Changed
- `GetComponent` now bypasses version resolution when the requirement is a fixed version (no range operators)
  - Returns directly without database lookup for exact versions
  - Covers operators from npm (`^~<>=*|`), pip (`~=!<>=`), maven/nuget (`[](),`), cargo, composer, and gems

## [0.2.0] - 2025-08-08
### Changed
- Refactored model constructors to use `*sqlx.DB` instead of `DBQueryContext`, removing `go-grpc-helper` dependency
- Simplified model constructors by removing logger and context parameters
- Renamed `LicenseID` field to `SPDX` in License model

## [0.1.1] - 2025-07-24
### Fixed
- Use `LicenseID` as ID from the DB

## [0.1.0] - 2025-07-24
### Added
- Initial release with shared models for SCANOSS services
- Added `AllUrlsModel` for URL lookups
- Added `ProjectModel` for project data access
- Added `VersionModel` for version data access
- Added `LicenseModel` for license data access
- Added `MineModel` for mine data access
- Added unified `Models` struct for accessing all models

[0.1.0]: https://github.com/scanoss/go-models/compare/v0.0.1...v0.1.0
[0.1.1]: https://github.com/scanoss/go-models/compare/v0.1.0...v0.1.1
[0.2.0]: https://github.com/scanoss/go-models/compare/v0.1.1...v0.2.0
[0.3.0]: https://github.com/scanoss/go-models/compare/v0.2.0...v0.3.0
[0.4.0]: https://github.com/scanoss/go-models/compare/v0.3.0...v0.4.0
[0.5.0]: https://github.com/scanoss/go-models/compare/v0.4.0...v0.5.0
[0.5.1]: https://github.com/scanoss/go-models/compare/v0.5.0...v0.5.1