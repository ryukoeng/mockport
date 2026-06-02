## Summary

## Test evidence

- [ ] Spec-first TDD evidence is included when behavior changed: spec/docs update, failing regression or contract test, implementation, and final verification.
- [ ] `/usr/local/go/bin/go test ./...`
- [ ] `/usr/local/go/bin/go vet ./...`
- [ ] `bash scripts/check-public-trust.sh`
- [ ] `bash scripts/check-distribution.sh`
- [ ] Docker or smoke test evidence if runtime behavior changed

## Public env safety

- [ ] No real provider secrets, production URLs, customer data, or unsanitized fixtures are included.
- [ ] Generated or documented credentials use Mockport fake values.

## Adapter changes

- [ ] Adapter spec, fixtures/manifests or SDK contract evidence, and known gaps were updated when the supported surface changed.
- [ ] Metadata/report coverage updated if adapter behavior changed.
- [ ] `mockport add <adapter>` and `mockport help <service>` still match adapter metadata.
- [ ] Unsupported behavior is documented or reported.
