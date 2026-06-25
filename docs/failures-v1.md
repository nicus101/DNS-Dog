# DNS-Dog Failures v1

DNS-Dog classifies failures as hard or transient.

Hard failures are configuration, setup, or permission problems that require
operator intervention before retrying.

Transient failures are runtime problems that may succeed later without changing
configuration.

## Hard failures

Hard failures cause `dns-dog run` and `dns-dog daemon` to exit non-zero.

| Failure | Classification |
| --- | --- |
| Invalid TOML | hard |
| Unsupported config format | hard |
| Missing required config | hard |
| Invalid duration or non-positive duration | hard |
| Invalid hook configuration | hard |
| No IP provider configured | hard |
| No credentials for configured reconciler | hard |
| Invalid credentials detected at startup | hard |
| Insufficient provider permissions detected at startup | hard |
| Configured reconciler cannot be initialized | hard |

## Transient failures

In `dns-dog daemon`, transient failures are logged and retried with backoff.

In `dns-dog run`, transient failures still cause a non-zero exit because there is
no later cycle in the same process.

| Failure | Classification |
| --- | --- |
| IP provider timeout | transient |
| IP provider temporary outage | transient |
| IP provider invalid response | transient |
| All IP providers fail during a cycle | transient |
| Reverse DNS lookup failure | transient |
| Provider timeout during reconciliation | transient |
| Provider rate limit during reconciliation | transient |
| Provider server error during reconciliation | transient |
| Zone refresh failure | transient |
| Hook command failure | transient |
| Hook timeout | transient |
| State file read/write I/O error | transient |

## Reconciliation edge cases

A missing configured DNS record is transient during normal operation. It may be a
temporary provider inconsistency or an operator-created record that has not
settled yet. Startup validation may classify the same condition as hard if it
proves the configured reconciler target is not accessible.

A malformed existing DNS record value is transient during normal operation. The
daemon retries so an operator or provider-side correction can recover without
changing DNS-Dog configuration.

Partial reconciliation is possible when DNS-Dog updates one record and then a
later record or zone refresh fails. The cycle fails, hooks do not run, and state
is not saved. The next cycle re-reads provider state and reconciles from the
current provider truth.

## Hook failure behavior

Hooks run in config order.

If a hook fails or times out, remaining hooks are skipped. The cycle fails and
state is not saved. In daemon mode, hooks are retried after the next successful
observation and reconciliation cycle.

DNS-Dog v1 does not support hooks that fire on failure. Use a separate
DNS-Dog instance, service manager policy, or external monitor/wrapper for failure
notification or recovery workflows.
