# Planejamento do Projeto

Este repositório contém a resolução do desafio para DevOps Engineer. O foco da implementação foi realizar "o básico bem feito", com uma arquitetura moderna, escalável e de fácil reprodutibilidade.

## Decisões Arquiteturais

### 1. Aplicações
- **App 1 (Python):** Utiliza FastAPI. Framework moderno, rápido e com documentação automática nativa, o que acelera o desenvolvimento de APIs e demonstra uso de ferramentas atualizadas no ecossistema Python.
- **App 2 (Go):** Utiliza a biblioteca padrão `net/http`. Demonstra conhecimento sólido dos fundamentos da linguagem sem depender de frameworks externos para tarefas simples.

### 2. Cache
- **Redis:** Uma instância centralizada para gerenciar o cache de ambas as aplicações. É o padrão de mercado para caching em memória.
- Foram implementados tempos de expiração (TTL) distintos para cada aplicação nas rotas que buscam o horário, cumprindo as regras do desafio (10s para Python, 1m para Go).

### 3. Observabilidade
- A monitoração local foca em simplicidade e efetividade:
  - **Prometheus:** Coleta as métricas expostas pelas aplicações (rota `/metrics`) e as métricas do cAdvisor.
  - **cAdvisor:** Agente oficial do Google para monitoramento de recursos de contêineres Docker, essencial para visualizar o consumo de CPU/Memória.
  - **Grafana:** Painel de visualização configurado para consolidar as métricas de infraestrutura e aplicação em um único lugar.

### 4. Facilidade de Execução (Docker Compose)
- Toda a stack é orquestrada pelo Docker Compose.
- Em vez de fazer o build local (que pode demorar e depender do poder computacional da máquina), o arquivo `docker-compose.yml` faz o pull das imagens já pré-construídas a partir de um Container Registry, agilizando o setup.

### 5. Integração Contínua (CI)
- **GitHub Actions:** O repositório possui uma pipeline que constrói e publica (push) as imagens no **GitHub Container Registry (GHCR)** automaticamente a cada alteração na branch principal, mantendo as imagens sempre atualizadas de forma transparente.
