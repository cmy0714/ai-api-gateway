# AI API Gateway

Modified version of [new-api](https://github.com/QuantumNous/new-api) — an LLM gateway and AI asset management platform.

**Corresponding source:** https://github.com/cmy0714/ai-api-gateway  
**Changes in this fork:** [MODIFICATIONS.md](./MODIFICATIONS.md)

[简体中文](./README.zh_CN.md)

## Quick Start

```bash
git clone https://github.com/cmy0714/ai-api-gateway.git
cd ai-api-gateway
docker-compose up -d
```

After deployment, open `http://localhost:3000`.

Build from source:

```bash
# backend
go build -o ai-api-gateway .

# frontend (web/default)
cd web/default && bun install && bun run build
```

See [docs.newapi.pro](https://docs.newapi.pro/en/docs/installation) for upstream deployment and configuration details.

## Fork Changes (Summary)

- Open payment integration (top-up, refund, wallet UI)
- Payment settings in system admin
- API key integration docs: [docs/api-key-integration.md](./docs/api-key-integration.md)
- Frontend home page and theme updates

Full list: [MODIFICATIONS.md](./MODIFICATIONS.md)

## Compliance Notes

> - Use only with lawful upstream API keys and applicable regulations.
> - When offering generative AI services to the public, comply with local filing, licensing, and content-safety requirements.

## Upstream

| Project | Link |
|---------|------|
| new-api (upstream) | https://github.com/QuantumNous/new-api |
| One API (original base) | https://github.com/songquanpeng/one-api |

## License

This project is licensed under the [GNU Affero General Public License v3.0 (AGPLv3)](./LICENSE).

Additional terms under AGPLv3 Section 7 apply. Modified versions must preserve the author attribution notice `Frontend design and development by New API contributors.` in the appropriate legal notices and in any prominent about, legal, footer, or attribution location presented by the user interface.

Modified versions that present a user interface must also preserve a visible link to the original project: https://github.com/QuantumNous/new-api

This fork is developed based on [One API](https://github.com/songquanpeng/one-api) (MIT License) via new-api.

For commercial licensing: [support@quantumnous.com](mailto:support@quantumnous.com)
