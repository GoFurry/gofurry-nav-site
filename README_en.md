# GoFurry

[中文说明](./README_zh.md)

GoFurry is a public-interest furry website and an open-source monorepo that contains the source code for its public site, data services, collectors, and operations backend.

This repository is organized by service. Each service can be developed and deployed independently, while sharing the same overall data model and infrastructure style.

## What This Repository Contains

- `gofurry-nav-frontend`: Vue frontend for the public navigation site
- `gofurry-nav-backend`: Go backend for navigation data APIs
- `gofurry-nav-collector`: Go collector for navigation-related data
- `gofurry-game-backend`: Go backend for game-related APIs
- `gofurry-game-collector`: Go collector for game-related data
- `gofurry-admin`: Operations backend with embedded Vue UI for daily maintenance
- `experimental`: experimental code, not part of production packaging
- `tools`: helper scripts and local tools, not part of production packaging

## Stack

- Go
- Fiber
- PostgreSQL
- Redis
- Vue
- Tailwind CSS

## Build

The root `build.bat` builds Linux `amd64` production artifacts into the root `build/` directory.

Build everything:

```bat
build.bat all
```

Build a single service:

```bat
build.bat gofurry-nav-backend
build.bat gofurry-nav-collector
build.bat gofurry-nav-frontend
build.bat gofurry-game-backend
build.bat gofurry-game-collector
build.bat gofurry-admin
```

Notes:

- Go binaries are built with production-oriented size reduction flags.
- `gofurry-admin` embeds its Vue frontend into the final binary.
- `experimental` and `tools` are intentionally excluded from packaging.

## Development

Each service is self-contained and should be started from its own directory.

Typical local workflow:

1. Enter the target service directory.
2. Install dependencies for that service.
3. Prepare your local configuration and database/Redis settings.
4. Run the service locally.

For frontend services:

- install dependencies with `npm`
- run the dev server from the frontend project directory

For Go services:

- use `go run . serve` or the service-specific startup command

## Deployment

Production deployment is expected to use private configuration files prepared by the deployer.

This repository does not ship production secrets. Do not commit:

- production database addresses
- Redis passwords
- JWT secrets
- TLS private keys
- production `server.yaml` or equivalent private config files

The repository `.gitignore` is configured to avoid committing common sensitive files, but it cannot erase anything already pushed in history. If a secret was ever committed, rotate it.

## Project Position

GoFurry is maintained as a public-interest project. This repository is opened so the site can be studied, improved, and extended more transparently.

The codebase is service-oriented instead of forcing everything into one runtime. That keeps deployment practical and allows different parts of the site to evolve at different speeds.

## Contributing

Issues and pull requests are welcome.

When contributing:

- keep changes scoped to the relevant service
- avoid committing local or production secrets
- preserve existing service boundaries unless there is a strong reason to change them

## License

See [LICENSE](./LICENSE).
