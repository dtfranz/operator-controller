# OCPSTRAT-3040: Adopt cluster-admin scope and deprecate ServiceAccount

Tracking checklist for the implementation work. Single shared feature branch, one PR.
Work is organized by file/code area ownership to maximize parallelism.

## Parallelism

Streams 1-3 start in parallel on day one. Stream 4 integrates after 1-3.

```
Engineer A:  Stream 1 (api/ + helm/)
Engineer B:  Stream 2 (authentication/ + action/ + labels/)
Engineer C:  Stream 3 (authorization/ + hack/tools/ + Makefile + go.mod)
             ────── all three in parallel ──────
Engineer D:  Stream 4 (appliers + controllers + features.go + main.go) ← after 1-3
             Stream 5 (contentmanager/ refactor) ← after 4
Anyone:      Stream 6 (docs) ← after 4
```

## Implementation

### Stream 1 — API + Helm (Engineer A, no blockers)

- [X] CRD schema change: `spec.serviceAccount` required → optional — **OPRUN-4630** `7f5a6843`
- [X] CRD upgrade safety validated — **OPRUN-4630** `7f5a6843`
- [X] Emit deprecation warning when `spec.serviceAccount` is set — **OPRUN-4630** `7f5a6843`
- [X] Grant operator-controller cluster-admin ClusterRole in helm RBAC templates — **OPRUN-4630** `7f5a6843`
- [X] Fix `Deprecated:` godoc format to follow Go/OpenShift conventions — **OPRUN-4630** `31cf46b7`
- [X] Add `omitzero` JSON tag to `serviceAccount` field — **OPRUN-4630** `31cf46b7`
- [X] Remove SA-based branching from `SyntheticUserRestConfigMapper` (`action/restconfig.go`) — **OPRUN-4630** `31cf46b7`
- [X] Fix deprecation warning log level and message in `ServiceAccountDeprecationWarning` — **OPRUN-4630** `31cf46b7`
- [X] Remove stale SA reference from `namespace` field godoc — **OPRUN-4630** `31cf46b7`
- [X] Add ValidatingAdmissionPolicy to emit kubectl deprecation warning when `serviceAccount` is set — **OPRUN-4630** `16b9c416`
- [X] Remove `serviceAccount` from all e2e feature files, add `namespace` step, restore incorrectly removed scenarios, clean up dead step code — **OPRUN-4630, OPRUN-4631, OPRUN-4632** `21336533`

### Stream 2 — Authentication + Action + Labels (Engineer B, no blockers)

- [ ] Remove `TokenGetter` (`authentication/tokengetter.go`) — **OPRUN-4630**
- [ ] Remove `TokenInjectingRoundTripper` (`authentication/tripper.go`) — **OPRUN-4630**
- [ ] Remove `authentication/synthetic.go` + tests — **OPRUN-4632**
- [ ] Remove `ServiceAccountRestConfigMapper` from `action/restconfig.go` — **OPRUN-4630**
- [ ] Remove `SyntheticUserRestConfigMapper` from `action/restconfig.go` + tests — **OPRUN-4632**
- [X] Remove SA annotations (`ServiceAccountNameKey`, `ServiceAccountNamespaceKey`) from `labels/labels.go` — **OPRUN-4630** `7f5a6843`

### Stream 3 — Authorization + Build tooling (Engineer C, no blockers)

- [X] Remove `authorization/` package entirely (~700 lines) — **OPRUN-4631**
- [X] Remove `k8s.io/kubernetes` dependency (require + replace in go.mod) — **OPRUN-4631**
- [X] Remove `hack/tools/k8smaintainer/` — **OPRUN-4631**
- [X] Remove `k8s-pin` Makefile target, simplify `verify` — **OPRUN-4631**

### Stream 4 — Integration: Appliers + Controllers + main.go (Engineer D, after streams 1-3)

- [X] Remove `ServiceAccountValidator` reconciliation step — **OPRUN-4630** `7f5a6843`
- [X] Remove `getUserInfo()` from Helm applier — **OPRUN-4630** `7f5a6843`
- [X] Remove `getUserInfo()` from Boxcutter applier — **OPRUN-4630** `7f5a6843`
- [ ] Remove pre-auth checks from Helm applier — **OPRUN-4631**
- [ ] Remove pre-auth checks from Boxcutter applier — **OPRUN-4631**
- [ ] Remove `PreflightPermissions` feature gate from `features.go` — **OPRUN-4631**
- [ ] Remove `SyntheticPermissions` feature gate from `features.go` — **OPRUN-4632**
- [ ] Clean up RBAC watches if only needed for pre-auth — **OPRUN-4631**
- [X] Rewire `main.go`: remove all SA/synthetic/preflight wiring, always use controller's own SA — **OPRUN-4630** `7f5a6843`
- [ ] Update all unit + e2e tests — **OPRUN-4630, OPRUN-4631, OPRUN-4632**

> **Ticket close gates:**
> - OPRUN-4630 closable when all streams 1, 2, 4 items marked with OPRUN-4630 are done
> - OPRUN-4631 closable when all streams 3, 4 items marked with OPRUN-4631 are done
> - OPRUN-4632 closable when all streams 2, 4 items marked with OPRUN-4632 are done

### Stream 5 — Global shared cache refactor (after stream 4)

- [ ] Refactor `contentmanager/` from per-CE caches to single global cache — **OPRUN-4633**
- [ ] Remove `createScopedClient` from `revision_engine_factory.go` — **OPRUN-4633**
- [ ] Simplify `RestConfigMapper` pattern — **OPRUN-4633**
- [ ] Refactor Boxcutter to shared client (if applicable) — **OPRUN-4633**
- [ ] Verify no regression in drift detection / managed content watching — **OPRUN-4633**
- [ ] Update unit + e2e tests — **OPRUN-4633**

> **Ticket close gate:** OPRUN-4633 closable when all stream 5 items are done

### Stream 6 — RBAC + VAP documentation (after stream 4)

- [ ] RBAC examples for common delegation scenarios — **OPRUN-4634**
- [ ] ValidatingAdmissionPolicy examples (catalog source, package, namespace restrictions) — **OPRUN-4634**
- [ ] Migration guidance from SA-based permission control — **OPRUN-4634**
- [ ] Update `serviceAccount` deprecation godoc in `api/v1/clusterextension_types.go` to reference RBAC + ValidatingAdmissionPolicy guidance once documented — **OPRUN-4634**

> **Ticket close gate:** OPRUN-4634 closable when all stream 6 items are done

## Post-merge

- [ ] Enhancement proposal in openshift/enhancements — **OPRUN-4629**
- [ ] API review for `spec.serviceAccount` deprecation — **OPRUN-4629**
- [ ] TP product docs — **OSDOCS-18848**

> **Ticket close gate:** OPRUN-4629 closable when EP is merged and API review is complete
