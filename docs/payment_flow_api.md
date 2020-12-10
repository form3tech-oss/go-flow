# Fake Payment API Flow

A sample payment flow to be demonstrated in the sample `go-flows` app. Not intended to directly match current or future Form3 flows, just a a similar example.

## Part 1 - create a payment submission resource

- `gin source` - POST to /v2/payments
- Validate & parse json - 400 if not valid
- Check token is authorized to post for this organisation
- `postgres sink` Attempt to persist to postgres 
    - Duplicate error - respond with `409 duplicate`
- `gin` respond with `204 no content`


## Part 2 - Validation

- `CDC source` Payment successfully persisted
- In parallel:
    - `sns sink` Notify subscribers & billing - SNS
    - `nats sink` Start generic payment API validation - (beneficiary validation api)
    - `nats sink` Notify Gateways for validation - SNS

### Part 2b Gateway Validation

* `nats source` payment submitted
* `nats sink` send validation results

## Part 3 - Validation results are in

* `nats source` Validation complete
* `postgres sink` - persist validation results

## Part 4 - Dispatch to the gateway

* `CDC source` - validation results persisted
* Check if we have both API & gateway validation, otherwise stop here.
* Something to do with ledger & limit checks
* More interesting flows - fraud checks, high value approval, etc.
* `nats sink` Send to gateway

## Part 4b - Gateway processing

* `nats source`
* `nats sink`

## Part 5 - Gateway response

* `nats source`
* `postgres sink` Persist response
* `sns sink` Notify subscribers

## Returns Flow

Slightly more interesting as we need to check the existing resource