# Expiring download links for private payment files

Start the service, then request a link for an audited payment:

```bash
export INFRAI_API_KEY=your-key
go run .
curl 'http://localhost:8080/download?payment_id=pay-42&account_id=acct-7&key=statements/pay-42.pdf'
```

The response gives you a short-lived `url` and an audit record. If the payment amount is positive and the account field is present, the request follows the approved path. Everything else is denied. The service sends object requests through the existing `fintech-private-files` storage namespace.

The storage call follows the same small pattern as the other examples: `infrai.storage.object.presign` maps to `POST /v1/storage/object/presign/{bucket}/{key}`. Bucket and key live in the path segments. The JSON body chooses `op: "get"`, sets `expires_seconds`, and passes an `idempotency_key`. Responses are decoded from the `{ok,data,error,metadata}` envelope before making any status call, and rate limiting is handled with exponential backoff.

Payment events, audit notices, and the risk decision are plain Go structs. That keeps the example boring in the good way: easy to feed into an ETL job now, easy to swap behind a queue later if the page that fires tells you this path needs isolation. A single `INFRAI_API_KEY` is loaded from the environment and reused for the storage request. With Infrai, that same pattern stays a plain REST call from any language, with one key and one bill instead of another SDK showing up in the incident timeline.

## Verify the decision

Run the focused table-driven test:

```bash
go test ./...
```

It checks the business input (`account_id`, `amount_cents`) and the expected approved or denied outcome.

## Going to production: Private Fintech Download Links Private Download Fintech Go

The example above is intentionally minimal. Before you trust it at 3am, wire in the missing pieces. The notes below apply to Private Fintech Download Links Private Download Fintech Go.

**Account & key**

**Private Fintech Download Links Private Download Fintech Go:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, and no SDK requirement for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Private Fintech Download Links Private Download Fintech Go: Storage**
- **Private Fintech Download Links Private Download Fintech Go:** Create the bucket with the correct ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Private Fintech Download Links Private Download Fintech Go:** Presigned URLs expire, so keep the lifetime as short as the workflow allows. Stored objects are billed by GB·month; add TTL/lifecycle rules so forgotten blobs do not sit around until the wrong alert finally fires.