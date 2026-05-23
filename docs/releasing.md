# Releasing VXT

## Preconditions

- clean working tree
- all intended release commits already present on `main`
- release notes prepared in `docs/releases/<version>.md`

## Verification

Run these commands from the repository root:

```bash
go test ./...
go list ./...
go doc github.com/vandordev/vxt
go doc github.com/vandordev/vxt/runtime
go doc github.com/vandordev/vxt/write
./scripts/release-smoke.sh
```

Before the public tag is resolvable remotely, `./scripts/release-smoke.sh` uses
a local `go mod edit -replace` directive to verify the consumer path against the
current checkout.

## Tagging

Create the release tag only after the verification sequence passes:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Post-Tag Checks

- verify the tagged revision still passes the same verification commands
- publish or paste the curated release notes from `docs/releases/v0.1.0.md`
- record any follow-up fixes separately instead of mutating the release checklist
