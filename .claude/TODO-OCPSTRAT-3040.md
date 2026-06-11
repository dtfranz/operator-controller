# OCPSTRAT-3040: Adopt cluster-admin scope and deprecate ServiceAccount

Tracking checklist for the implementation work. Single shared feature branch, one PR.

## Implementation

### Commit 1 — Make spec.serviceAccount optional, grant cluster-admin RBAC (OPRUN-4630)
- [ ] CRD schema change: `spec.serviceAccount` required → optional
- [ ] CRD upgrade safety validated
- [ ] Emit deprecation warning when `spec.serviceAccount` is set
- [ ] Grant operator-controller cluster-admin ClusterRole in helm RBAC templates
- [ ] Update unit + e2e tests

### Commit 2 — Remove SA-scoped client infra (OPRUN-4630)
- [ ] Remove `TokenGetter` (`authentication/tokengetter.go`)
- [ ] Remove `TokenInjectingRoundTripper` (`authentication/tripper.go`)
- [ ] Remove `ServiceAccountRestConfigMapper` (`action/restconfig.go`)
- [ ] Remove `ServiceAccountValidator` reconciliation step
- [ ] Remove `getUserInfo()` from Helm and Boxcutter appliers
- [ ] Remove SA annotations (`ServiceAccountNameKey`, `ServiceAccountNamespaceKey`)
- [ ] Rewire `main.go` to always use controller's own SA
- [ ] Update unit + e2e tests

### Commit 3 — Remove SyntheticPermissions (OPRUN-4632)
- [ ] Remove `SyntheticPermissions` feature gate from `features.go`
- [ ] Remove `authentication/synthetic.go` + tests
- [ ] Remove `SyntheticUserRestConfigMapper` from `action/restconfig.go` + tests
- [ ] Remove wiring from `main.go`
- [ ] Update unit + e2e tests

### Commit 4 — Remove PreflightPermissions (OPRUN-4631)
- [ ] Remove `PreflightPermissions` feature gate from `features.go`
- [ ] Remove `authorization/` package entirely (~700 lines)
- [ ] Remove pre-auth checks from Helm applier
- [ ] Remove pre-auth checks from Boxcutter applier
- [ ] Remove `k8s.io/kubernetes` dependency (require + replace)
- [ ] Remove `hack/tools/k8smaintainer/`
- [ ] Remove `k8s-pin` Makefile target, simplify `verify`
- [ ] Clean up RBAC watches if only needed for pre-auth
- [ ] Remove wiring from `main.go`
- [ ] Update unit + e2e tests

### Commit 5 — Consolidate to global shared cache (OPRUN-4633)
- [ ] Refactor `contentmanager/` from per-CE caches to single global cache
- [ ] Remove `createScopedClient` from `revision_engine_factory.go`
- [ ] Simplify `RestConfigMapper` pattern
- [ ] Refactor Boxcutter to shared client (if applicable)
- [ ] Verify no regression in drift detection / managed content watching
- [ ] Update unit + e2e tests

### Commit 6 — RBAC + VAP documentation (OPRUN-4634)
- [ ] RBAC examples for common delegation scenarios
- [ ] ValidatingAdmissionPolicy examples (catalog source, package, namespace restrictions)
- [ ] Migration guidance from SA-based permission control

## Post-merge

- [ ] Enhancement proposal in openshift/enhancements (OPRUN-4629)
- [ ] API review for `spec.serviceAccount` deprecation
- [ ] TP product docs (OSDOCS-18848)

## Parallelism

```
Engineer A:  commits 1 → 2 ·············→ 5
Engineer B:       commit 3 → 4
Anyone:                                    6
```
