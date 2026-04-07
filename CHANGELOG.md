# Changelog

## [0.1.3] - 2026-04-07

**Target any repository directly.** Generate credits for the current repo or pass a repository path explicitly, including GIF export for another repo.

```bash
gitcredits /path/to/repo
gitcredits /path/to/repo --output credits.gif
```

### Added

- support a target repository directory as a positional argument
- add tests for target directory parsing and invalid path handling

### Changed

- parse the repository directory as a proper positional argument instead of relying on fragile argument order
- access repository data by explicit path instead of changing the global working directory
- return errors from helpers and let `main()` handle exit behavior
- pass the target repository path through GIF generation so `--output` works for another repo

## [0.1.2] - 2026-03-21

**New Spider-Man theme with glitch effects and radial web transitions.**

```bash
gitcredits --theme spiderman
```

### Added

- Spider-Man theme: RGB shift, glitch text distortion, radial web-shooting transitions
- Spider-Man character titles for contributors
- Notable commits and full stats cards for Spider-Man theme

### Fixed

- VHS now runs in current working directory (correct repo info in GIF output)
- Updated demo GIFs with correct project data

## [0.1.1] - 2026-02-28

### Added

- Matrix theme with digital rain effect
- GIF export support

## [0.1.0] - 2026-02-27

- Initial release
- Star Wars style rolling credits
- Big text title rendering
