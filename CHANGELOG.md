# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2025-03-18

### Added
- Support for enabling log monitoring with OAuth authentication
- Improved error handling for connector operations
- Added namespaceTag field to Azure EventHub configuration

### Changed
- Modified update process to exclude sensitive fields in auth_params
- Enhanced comparison logic for nested structures
- Updated documentation with more detailed examples

### Fixed
- Fixed issue with updating connectors that have sensitive fields
- Resolved internal error when enabling log monitoring

## [0.3.0] - 2025-03-18

### Added
- Enhanced error handling for deleted connectors
- Support for recreating connectors that were deleted outside of Terraform
- Added status, enabled state, last activity, and outpost ID fields

### Changed
- Improved retry logic with exponential backoff
- Enhanced deep comparison for nested structures

### Fixed
- Fixed issue with connector recreation after external deletion

## [0.2.0] - 2025-03-18

### Added
- State comparison before making changes
- Update support instead of recreating resources
- Enhanced schema with additional fields

### Changed
- Improved error handling and logging
- Better validation of input parameters

### Fixed
- Various bug fixes and stability improvements

## [0.1.0] - 2025-03-17

### Added
- Initial release of the Terraform Provider for Wiz
- Support for creating, reading, and deleting Wiz connectors
- Support for testing connector configurations
- Basic documentation and examples
