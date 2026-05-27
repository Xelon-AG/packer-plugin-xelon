
## v0.3.0 (2026-05-27)
### Code Refactoring
* extract AccessConfig from builder
### Documentation
* prepare web-docs for integrations listing
### Features
* add datasource to fetch template data
* add datasource to fetch network data
### Maintaining
* fix test for access config logic
* **dependabot**: reduce noise with all hashicorp deps
* **dependabot**: add tools for auto-checking package updates from dependabot
* **deps**: bump github.com/hashicorp/packer-plugin-sdk in /tools
* **deps**: bump github.com/hashicorp/packer-plugin-sdk
* **deps**: bump goreleaser/goreleaser-action from 7.2.1 to 7.2.2
* **deps**: bump github.com/golangci/golangci-lint/v2 in /tools
* **deps**: bump github.com/hashicorp/packer-plugin-sdk in /tools
* **deps**: bump github.com/hashicorp/packer-plugin-sdk

## v0.2.0 (2026-05-18)
### Bug Fixes
* set default template description
* **deps**: upgrade crypto to v0.51.0
### Code Refactoring
* add validation for AccessConfig
* split and extract config object from builder
### Features
* provide template-related properties via generated_data
### Maintaining
* fix golangci linter issues
* **actions**: sign release binaries with gpg key
* **actions**: refactor test github action for pull requests

## v0.1.1 (2026-05-11)
### Bug Fixes
* set pre-release ldflag correctly

## v0.1.0 (2026-05-11)
### Bug Fixes
* query correct ip address from device
### Documentation
* remove scaffolding docs
### Maintaining
* **actions**: prepare release gh action
* **deps**: upgrade goreleaser/goreleaser-action to v7.2.1
* **deps**: upgrade crazy-max/ghaction-import-gpg to v7.0.0
* **deps**: upgrade actions/setup-go to v6.4.0
* **deps**: upgrade actions/checkout to v6.0.2
