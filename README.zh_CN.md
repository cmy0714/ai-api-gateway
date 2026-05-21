# AI API Gateway

基于 [new-api](https://github.com/QuantumNous/new-api) 的修改版本，提供大模型网关与 AI 资产管理能力。

**对应源码：** https://github.com/cmy0714/ai-api-gateway  
**本仓库变更说明：** [MODIFICATIONS.md](./MODIFICATIONS.md)

[English](./README.md)

## 快速开始

```bash
git clone https://github.com/cmy0714/ai-api-gateway.git
cd ai-api-gateway
docker-compose up -d
```

部署完成后访问 `http://localhost:3000`。

从源码构建：

```bash
# 后端
go build -o ai-api-gateway .

# 前端（web/default）
cd web/default && bun install && bun run build
```

更多部署与环境变量说明请参考上游文档：[docs.newapi.pro](https://docs.newapi.pro/zh/docs/installation)

## 本 Fork 主要变更

- 开放支付集成（充值、退款、钱包 UI）
- 系统设置中的支付配置
- API Key 对接文档：[docs/api-key-integration.md](./docs/api-key-integration.md)
- 首页与主题样式调整

完整列表见 [MODIFICATIONS.md](./MODIFICATIONS.md)

## 合规说明

> - 请合法取得上游 API Key，并遵守上游服务条款及适用法律法规。
> - 面向公众提供生成式 AI 服务时，请自行完成所在地区要求的备案、许可与内容安全等义务。

## 上游项目

| 项目 | 链接 |
|------|------|
| new-api（上游） | https://github.com/QuantumNous/new-api |
| One API（原始基础） | https://github.com/songquanpeng/one-api |

## 许可证

本项目采用 [GNU Affero 通用公共许可证 v3.0 (AGPLv3)](./LICENSE) 授权。

根据 AGPLv3 第 7 条附加条款，修改版须在适当的法律声明以及界面中显眼的 about、legal、footer 或 attribution 位置保留以下作者归属说明：

`Frontend design and development by New API contributors.`

提供用户界面的修改版还须在显眼位置保留指向上游项目的可见链接：https://github.com/QuantumNous/new-api

本 Fork 基于 [One API](https://github.com/songquanpeng/one-api)（MIT 许可证）经由 new-api 二次开发。

商业授权请联系：[support@quantumnous.com](mailto:support@quantumnous.com)
