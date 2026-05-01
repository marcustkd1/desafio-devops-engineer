# Arquitetura da Solução

O diagrama a seguir descreve a topologia da infraestrutura proposta e o fluxo da integração contínua. 
Pode ser visualizado nativamente no GitHub ou colado em ferramentas compatíveis com Mermaid (ex: Miro, Notion).

```mermaid
flowchart TD
    %% Nós do usuário
    User((Usuário))

    subgraph CI/CD [Pipeline GitHub Actions]
        CodePush[Commit no Repositório] --> BuildImages(Build das Imagens Docker)
        BuildImages --> PushGHCR[(GitHub Container Registry)]
    end

    subgraph Infraestrutura Local [Docker Compose]
        %% Aplicações
        AppPython(App 1 - Python FastAPI\nPorta: 8000)
        AppGo(App 2 - Go net/http\nPorta: 8080)
        
        %% Componentes auxiliares
        Redis[(Redis Cache)]
        
        %% Observabilidade
        Prometheus(Prometheus\nPorta: 9090)
        cAdvisor(cAdvisor)
        Grafana(Grafana\nPorta: 3000)
    end

    %% Relações de Tráfego do Usuário
    User -- "Acessa APIs" --> AppPython
    User -- "Acessa APIs" --> AppGo
    User -- "Visualiza Dashboards" --> Grafana
    
    %% Relações de Banco
    AppPython -- "Cache TTL 10s" --> Redis
    AppGo -- "Cache TTL 60s" --> Redis
    
    %% Relações de Observabilidade
    Prometheus -. "Scrape /metrics" .-> AppPython
    Prometheus -. "Scrape /metrics" .-> AppGo
    Prometheus -. "Scrape métricas docker" .-> cAdvisor
    Grafana -. "Consulta de dados" .-> Prometheus
    
    %% Relacao de Deploy
    PushGHCR -. "Pull Images" .-> AppPython
    PushGHCR -. "Pull Images" .-> AppGo
```
