# Desafio Técnico DevOps

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
docker compose up -d
```
> O Docker fará o pull das imagens mais recentes diretamente do GHCR e subirá todos os serviços.

### O que estará rodando?

| Serviço | URL de Acesso Local | Descrição |
| --- | --- | --- |
| **App Python** | [http://localhost:8000](http://localhost:8000) | API em FastAPI (Cache de 10s). Rotas: [Texto](http://localhost:8000/), [Hora](http://localhost:8000/hora), [Health](http://localhost:8000/health), [Métricas](http://localhost:8000/metrics) |
| **App Go** | [http://localhost:8080](http://localhost:8080) | API em Go nativo (Cache de 60s). Rotas: [Texto](http://localhost:8080/), [Hora](http://localhost:8080/hora), [Health](http://localhost:8080/health), [Métricas](http://localhost:8080/metrics) |
| **Grafana** | [http://localhost:3000](http://localhost:3000) | Painel de visualização de métricas (Login/Senha Padrão: `admin`). |
| **Prometheus**| [http://localhost:9090](http://localhost:9090) | Interface principal do coletor de métricas. |
| **Node Exporter** | [http://localhost:9100](http://localhost:9100) | Coletor de métricas da infraestrutura (CPU, RAM, Disco). |
| **Redis** | `localhost:6379` | Banco de dados em memória atuando como cache (Acesso via CLI ou Client, sem interface web). |
