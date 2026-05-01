# Desafio Técnico DevOps

Bem-vindo(a) ao repositório do desafio técnico para DevOps.

Este repositório contém a resolução completa do teste, contemplando duas aplicações (Python/FastAPI e Go), sistema de cache (Redis), observabilidade (Prometheus + Grafana + cAdvisor) e integração contínua via GitHub Actions.

## Documentação do Projeto

O detalhamento da arquitetura, os diagramas e os pontos de melhoria estão separados nos seguintes documentos:
- [**Planejamento e Arquitetura**](./planejamento.md)
- [**Diagrama da Solução**](./diagrama.md)
- [**Pontos de Melhoria**](./melhorias.md)

## Como Executar (Fácil e Rápido)

Para facilitar a execução, as imagens Docker de ambas as aplicações são geradas automaticamente na pipeline de CI e armazenadas no **GitHub Container Registry (GHCR)**. Sendo assim, você não precisa "buildar" as imagens localmente, economizando tempo e recursos da sua máquina.

### Pré-requisitos
- Docker e Docker Compose instalados.

### Passos
1. Clone este repositório:
```bash
git clone https://github.com/marcustkd1/desafio-devops-engineer.git
cd desafio-devops-engineer
```

2. Suba a infraestrutura completa:
```bash
docker-compose up -d
```
> O Docker fará o pull das imagens mais recentes diretamente do GHCR e subirá todos os serviços.

### O que estará rodando?

| Serviço | Porta | Descrição |
| --- | --- | --- |
| **App Python** | `8000` | Aplicação FastAPI com cache de 10s. Rotas: `/` e `/hora`. |
| **App Go** | `8080` | Aplicação Go nativa com cache de 60s. Rotas: `/` e `/hora`. |
| **Redis** | `6379` | Banco de dados em memória atuando como cache. |
| **Grafana** | `3000` | Painel de métricas. (Login Padrão: `admin` / Pode ser pulado na primeira tela). |
| **Prometheus** | `9090` | Interface do coletor de métricas das aplicações. |
| **cAdvisor** | `8081` | Interface com métricas de infraestrutura dos containers Docker. |

> Ambas as aplicações possuem também as rotas `/health` e `/metrics` ativas.
