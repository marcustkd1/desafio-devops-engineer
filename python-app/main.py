import os
import redis
from datetime import datetime
from fastapi import FastAPI
from prometheus_fastapi_instrumentator import Instrumentator

app = FastAPI(title="App 1 - Python")

# Configuração do Redis
REDIS_HOST = os.getenv("REDIS_HOST", "localhost")
REDIS_PORT = int(os.getenv("REDIS_PORT", 6379))
redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, db=0, decode_responses=True)

# Instrumentação para o Prometheus na rota /metrics
Instrumentator().instrument(app).expose(app)

@app.get("/")
def read_root():
    return {"message": "Olá! Esta é a aplicação em Python respondendo com um texto fixo."}

@app.get("/hora")
def get_time():
    cache_key = "python_app_time"
    cached_time = redis_client.get(cache_key)
    
    if cached_time:
        ttl = redis_client.ttl(cache_key)
        return {"source": "cache", "time": cached_time, "cache_key": cache_key, "ttl": ttl}
    
    current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    # Cache por 10 segundos
    redis_client.setex(cache_key, 10, current_time)
    
    return {"source": "server", "time": current_time, "cache_key": cache_key, "ttl": 10}

@app.get("/health")
def health_check():
    # Verifica conexão com Redis de forma básica
    try:
        redis_client.ping()
        redis_status = "ok"
    except redis.ConnectionError:
        redis_status = "error"
        
    return {
        "status": "up" if redis_status == "ok" else "degraded",
        "redis": redis_status
    }
