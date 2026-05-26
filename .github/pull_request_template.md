## Summary

## Test evidence

- [ ] `/usr/local/go/bin/go test ./...`
- [ ] `/usr/local/go/bin/go vet ./...`
- [ ] `bash scripts/check-public-trust.sh`
- [ ] `bash scripts/check-distribution.sh`
- [ ] Docker or smoke test evidence if runtime behavior changed

## Public env safety

- [ ] No real provider secrets, production URLs, customer data, or unsanitized fixtures are included.
- [ ] Generated or documented credentials use Mockport fake values.

## Adapter changes

- [ ] Metadata/report coverage updated if adapter behavior changed.
- [ ] Unsupported behavior is documented or reported.
