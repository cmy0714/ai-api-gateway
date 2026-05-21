# Modifications Notice (AGPLv3 Section 7)

This repository contains a modified version of [new-api](https://github.com/QuantumNous/new-api).

## Upstream

- Project: new-api
- Upstream repository: https://github.com/QuantumNous/new-api
- Upstream license: GNU Affero General Public License v3.0 (AGPL-3.0)
- Base commit: `2d1ca153` (main at time of fork)

## License

This modified version is distributed under the same **AGPL-3.0** license as upstream.
See [LICENSE](./LICENSE) and [NOTICE](./NOTICE) for full terms, including Section 7
attribution requirements.

## Corresponding Source

The complete corresponding source code for this modified version is available at:

https://github.com/cmy0714/new-api

When this software is offered as a network service, users interacting with
it remotely may obtain the corresponding source from the repository above.

## Summary of Changes

The following changes were made on top of upstream new-api:

### Open payment integration

- Added open payment client, signing, and configuration (`service/paymentopen/`, `setting/payment_open.go`)
- Added open payment top-up, webhook, refund, and status APIs (`controller/topup_open_payment.go`, `model/topup_refund.go`)
- Extended wallet UI with QR payment dialog, refund dialog, and polling hooks
- Added payment settings section in system settings

### Wallet and billing

- Extended wallet API, billing history, recharge flow, and payment method handling
- Updated top-up controller and model for open payment compatibility

### API keys

- Extended API key management API and types
- Added API key integration documentation (`docs/api-key-integration.md`)

### Frontend and branding

- Updated default theme home page sections (hero, features, CTA, stats, how-it-works)
- Updated auth layout, theme CSS, and site metadata handling
- Removed some default system info / site settings fields from admin UI

### Build and deployment

- Adjusted Dockerfile, makefile, release workflow, and embedded frontend assets

## Attribution Preserved

Modified versions retain the upstream attribution required by AGPLv3 Section 7 and the
project NOTICE file, including:

- Source file copyright headers
- Footer link to https://github.com/QuantumNous/new-api
- New API contributor attribution in the user interface

## Commercial Licensing

For commercial licensing that does not require AGPL obligations, contact:
support@quantumnous.com
