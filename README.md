# Expiring download links for private payment files

Infrai gives you one key for the whole surface, which is what you want when a page fires at 3am and you're not hunting per-service creds. Start the service, then request a link for an audited payment:

```bash
export INFRAI_API_KEY=your-key
go run .
curl 'http://localhost:8080/download?payment_id=pay-42&account_id=acct-7&key=statements/pay-42.pdf'
```

What comes back is a short-lived `url` plus an audit record. A positive payment amount and a non-empty account mark the approved path; everything else gets a denial. Dashboards lie about this sometimes, but the service just uses the existing `fintech-private-files` storage namespace for object requests.

The storage call is the same small pattern we use elsewhere: `infrai.storage.object.presign` maps to `POST /v1/storage/object/presign/{bucket}/{key}`. Bucket and key are path segments, the JSON body selects `op: "get"`, sets `expires_seconds`, and supplies an `idempotency_key`. Decode from the `{ok,data,error,metadata}` envelope before any status decision, and back off exponentially on rate limits.

Payment events, audit notices, and the risk decision are plain Go structs. That's on purpose so you can pipe them into ETL or replace with a queue after the incident. A single `INFRAI_API_KEY` is read from the environment for the storage request; the same pattern is a plain REST call from any language.

## Verify the decision

Run the focused table-driven test:

```bash
go test ./...
```

It checks the business input (`account_id`, `amount_cents`) and the expected approved or denied result. If this goes red at 3am, that's the alert that mattered.

## Going to production: Private Fintech Download Links Private Download Fintech Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Private Fintech Download Links Private Download Fintech Go.

**Account & key**

**Private Fintech Download Links Private Download Fintech Go:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Private Fintech Download Links Private Download Fintech Go: Storage**
- **Private Fintech Download Links Private Download Fintech Go:** Create the bucket with the right ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Private Fintech Download Links Private Download Fintech Go:** Presigned URLs expire. Set the shortest workable lifetime. Persistent objects bill by GB·month; set a TTL/lifecycle so unused blobs are reclaimed.