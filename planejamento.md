# Planejamento do Projeto

Este arquivo contém a minha linha de raciocínio para a resolução do desafio para DevOps Engineer. O foco da implementação foi realizar "o básico bem feito", com uma arquitetura moderna, com telemetria básica e de fácil reprodutibilidade.

## Decisões Arquiteturais

### 1. Aplicações e Instrumentação
- **App 1 (Python):** Utiliza FastAPI. Framework moderno, rápido e com documentação automática nativa e integração direta com o Prometheus.
- **App 2 (Go):** Utiliza a biblioteca padrão `net/http` e o módulo oficial `client_golang`. Para a telemetria, foi desenvolvido um middleware customizado capaz de rastrear em tempo real o Volume Absoluto de Requisições, a distribuição de Status Codes e a Latência de resposta (Histogramas).

### 2. Cache Estruturado
- **Redis:** Uma instância local atua como repositório rápido na memória, mitigando a necessidade de as aplicações processarem repetidamente dados estáticos ou dinâmicos de baixa volatilidade.
- Foram implementados tempos de expiração (TTL) distintos na rota de validação, provando o controle fino do motor de cache (10s para FastAPI e 60s para Go). A API responde dinamicamente se o dado foi processado pelo "server" ou devolvido pelo "cache".

### 3. Observabilidade e Telemetria
Ao invés de apenas confirmar que os contêineres estão em execução, a stack entrega dados críticos e comparativos:
- **Prometheus:** Coletor configurado para fazer *scrape* periódico (5s) das aplicações e do Host.
- **Node Exporter:** Fornece uma visão aprofundada a nível de sistema (CPU, RAM, Discos, I/O) do host.
- **Grafana e Autoprovisionamento:** Server já pré-configurado, com dois *Dashboards* default (Infra e Apps):
    - **Visão dos Apps:** Gráficos comparativos (Python vs Go) focando no **Volume de Chamadas**, **Tempos de Resposta** e **Erros (HTTP 5xx / 4xx)**.
    - **Visão da Infra:** Um detalhamento massivo da saúde do sistema, processado nativamente via Node Exporter.

### 4. Orquestração e Reprodutibilidade (Docker Compose)
- A premissa do desafio era a simplicidade. Toda a infraestrutura sobe através de uma única instrução (`docker compose up -d`).
- **Eficiência Local:** Ao invés de o projeto compilar as aplicações via `build`, o Compose faz apenas o *Pull* das Imagens pré-compiladas pela CI e armazenadas no GitHub Container Registry.

### 5. Integração Contínua Inteligente (GitHub Actions)
- **Multi-stage Build & GHCR:** Pipeline estruturada que constrói imagens otimizadas em Alpine Linux e as envia de forma segura para o GitHub Container Registry.
- **Filtro de Gatilho (Paths):** A pipeline foi refinada para não desperdiçar minutos e recursos computacionais do runner free do Github caso o *commit* atualize apenas documentações (`*.md`) ou arquivos de configuração de infraestrutura (`.yml`, `.yaml`, `.md`). Os builds só ocorrem caso haja mutação na base de código do Python ou do Go.
