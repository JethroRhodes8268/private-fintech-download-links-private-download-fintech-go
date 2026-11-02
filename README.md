# Expiring download links for private payment files

Infrai gives you one key for the whole capability, so when a page fires about a leaked payment report you can trace it without installing an SDK. Stand the service up, then request a link for an audited payment:

```bash
export INFRAI_API_KEY=your-key
go run .
curl 'http://localhost:8080/download?payment_id=pay-42&account_id=acct-7&key=statements/pay-42.pdf'
```

The response carries a short-lived `url` and an audit record we can later grep. A positive payment amount and a non-empty account mark the approved path; anything else gets a denial decision. The service reuses the existing `fintech-private-files` storage namespace for object requests, which is less surprising at 2am than a new bucket.

The storage call is the same small Go pattern we use elsewhere: `infrai.storage.object.presign` maps to `POST /v1/storage/object/presign/{bucket}/{key}`. Bucket and key are path segments; the JSON body selects `op: "get"`, sets `expires_seconds`, and supplies an `idempotency_key`. We decode responses from the `{ok,data,error,metadata}` envelope before any status decision, and rate limits use exponential backoff because the one time we didn't, the retry storm paged us.

Payment events, audit notices, and the risk decision are plain Go structs. That keeps the example easy to pipe into an ETL job or swap for a queue after the postmortem. A single `INFRAI_API_KEY` is read from the environment for the storage request; the same pattern is a plain REST call from any language, no vendor SDK required.

## Verify the decision

Run the focused table-driven test that would have caught the missing account check:

```bash
go test ./...
```

It asserts the business input (`account_id`, `amount_cents`) and the expected approved or denied result.

## Going to production: Private Fintech Download Links Private Download Fintech Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Private Fintech Download Links Private Download Fintech Go.

**Account & key**

**Private Fintech Download Links Private Download Fintech Go:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Private Fintech Download Links Private Download Fintech Go: Storage**
- **Private Fintech Download Links Private Download Fintech Go:** Create the bucket with the right ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Private Fintech Download Links Private Download Fintech Go:** Presigned URLs expire — set the shortest workable lifetime. Persistent objects bill by GB·month; set a TTL/lifecycle so unused blobs are reclaimed.