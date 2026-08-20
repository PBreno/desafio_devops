# 🚀 Desafio DevOps — Projeto Korp

Este repositório contém a solução para o desafio DevOps da **Korp**, contemplando a criação de um serviço HTTP em **Golang**, infraestrutura de containers com **Docker e NGINX**, monitoramento com **Prometheus e Grafana** e automação completa do provisionamento com **Ansible**.

---

## 📋 Sumário

- [Arquitetura do Projeto](#-arquitetura-do-projeto)
- [Estrutura de Pastas](#-estrutura-de-pastas)
- [Pré-requisitos](#-pré-requisitos)
- [Como Executar](#-como-executar)
- [Endpoints e Acessos](#-endpoints-e-acessos)
- [Monitoramento e Métricas](#-monitoramento-e-métricas)
- [Autor](#-autor)

---

## 🏗 Arquitetura do Projeto

A arquitetura é baseada em containers executando em uma **rede bridge isolada do Docker**.

O **NGINX** atua como ponto de entrada da aplicação e como **Proxy Reverso**, encaminhando as requisições para o serviço desenvolvido em Go.

O **Prometheus** coleta as métricas expostas pela aplicação, enquanto o **Grafana** é responsável pela visualização e acompanhamento dos indicadores.

```text
                         +----------------------+
                         |      Host Linux      |
                         |                      |
                         |   Requisição HTTP    |
                         |          :80         |
                         |           │          |
                         |           ▼          |
                         |   +--------------+   |
                         |   |    NGINX     |   |
                         |   | Proxy Reverso |   |
                         |   +------+-------+   |
                         |          │           |
                         |     Rede Bridge      |
                         |          │           |
                         |   +------+-------+   |
                         |   |  Golang App  |   |
                         |   |    :8080     |   |
                         |   +------+-------+   |
                         |          │           |
                         |      +---+---+       |
                         |      │       │       |
                         |  +---▼---+ +--▼---+  |
                         |  |Prom.  | |Graf. |  |
                         |  | :9090 | | :3000|  |
                         |  +-------+ +------+  |
                         +----------------------+
```

### Fluxo da requisição

```text
Cliente
   │
   ▼
NGINX :80
   │
   ▼
Golang :8080
   │
   ├── /projeto-korp
          │
          ▼
      Prometheus
          │
          ▼
       Grafana
```

> **Nota:** O container da aplicação Golang não é exposto diretamente ao host. A aplicação fica acessível externamente através do NGINX, reduzindo a superfície de exposição do serviço.

---

## 📁 Estrutura de Pastas

```text
desafio_devops/
├── ansible/                                  # Automação de infraestrutura
│   ├── ansible.cfg
│   ├── inventory.ini
│   └── playbook.yml
│
├── docker/                                  # Infraestrutura de containers
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── nginx/
│       └── http-server-projeto-korp.conf
│
├── monitoring/                              # Configurações de observabilidade
│   ├── grafana/
│   │   ├── dashboards/
│   │   │   └── http-server-projeto-korp-dashboard.json
│   │   └── provisioning/
│   │       ├── dashboards/
│   │       │   └── dashboards.yml
│   │       └── datasources/
│   │           └── datasource.yml
│   │
│   └── prometheus/
│       └── prometheus.yml
│
├── src/                                     # Código-fonte da aplicação
│   ├── go.mod
│   ├── go.sum
│   └── main.go
│
└── README.md
```

---

## 🛠 Pré-requisitos

O ambiente onde o playbook será executado precisa ter:

| Requisito | Descrição |
|---|---|
| **Sistema Operacional** | Linux — Ubuntu 22.04/24.04 ou Debian |
| **Ansible** | Instalado na máquina responsável por executar a automação |
| **Privilégios** | Acesso `sudo`/`root` para instalação de pacotes e configuração do ambiente |
| **Acesso à Internet** | Necessário para instalação dos pacotes e obtenção das imagens Docker |

> O Docker não precisa estar previamente instalado, pois sua instalação é realizada automaticamente pelo playbook.

---

## ⚡ Como Executar

Todo o ambiente é provisionado de forma **idempotente** através do Ansible.

O playbook é responsável por:

1. Instalar os pré-requisitos;
2. Instalar o Docker;
3. Copiar os arquivos do projeto;
4. Construir a imagem da aplicação;
5. Criar e configurar os containers;
6. Subir NGINX, aplicação, Prometheus e Grafana;
7. Validar o funcionamento do serviço.

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/desafio_devops.git
cd desafio_devops/ansible
```

### 2. Execute o playbook

```bash
ansible-playbook playbook.yml
```

Ao final da execução, o playbook deverá validar o funcionamento do serviço.

Exemplo de resposta:

```json
{
  "nome": "Projeto Korp",
  "horario": "2024-05-20T15:30:00.123Z"
}
```

### Execução com inventário específico

Caso seja necessário informar explicitamente o inventário:

```bash
ansible-playbook -i inventory.ini playbook.yml
```

---

## 🌐 Endpoints e Acessos

Após a execução do playbook, os seguintes serviços estarão disponíveis no host:

| Serviço | URL | Credenciais | Descrição |
|---|---|---|---|
| **Aplicação via NGINX** | `http://localhost/projeto-korp` | — | Retorna JSON contendo o nome do projeto e o horário atual em UTC |
| **Prometheus UI** | `http://localhost:9090` | — | Interface de consulta e monitoramento do Prometheus |
| **Grafana** | `http://localhost:3000` | `admin / admin` | Dashboard de monitoramento da aplicação |

### Exemplo — Healthcheck

```bash
curl http://localhost/health
```

Resposta:

```json
{
  "status": "ok"
}
```

### Exemplo — Aplicação

```bash
curl http://localhost/projeto-korp
```

Resposta:

```json
{
  "nome": "Projeto Korp",
  "horario": "2024-05-20T15:30:00.123Z"
}
```

---

## 📊 Monitoramento e Métricas

A aplicação disponibiliza métricas compatíveis com o **Prometheus**.

Foram implementadas duas métricas principais:

### 1. Disponibilidade do serviço

```text
http_server_up
```

Tipo:

```text
Gauge
```

Valores:

- `1` → serviço disponível
- `0` → serviço indisponível


### 2. Volume de requisições

```text
http_requests_total
```

Tipo:

```text
Counter
```

A métrica é segmentada por:

- Endpoint
- Método HTTP

Isso permite acompanhar o volume de tráfego e identificar quais endpoints estão sendo mais utilizados.

---

## 📈 Grafana

O Grafana é provisionado automaticamente através dos arquivos presentes em:

```text
monitoring/grafana/provisioning/
```

O ambiente é configurado sem necessidade de intervenção manual.

Ao acessar:

```text
http://localhost:3000
```

o datasource do Prometheus já estará configurado e o dashboard:

```text
Projeto Korp - Monitoramento
```

estará disponível.

O dashboard apresenta informações relacionadas a:

- Disponibilidade da aplicação;
- Volume de requisições;
- Tráfego por endpoint;
- Métodos HTTP utilizados.

---

### 🔄 Automação com Ansible

O provisionamento da infraestrutura é realizado através do Ansible.

O playbook foi desenvolvido para ser **idempotente**, permitindo sua execução repetida sem gerar configurações inconsistentes.

Em vez de depender de bibliotecas Python específicas, o playbook utiliza comandos nativos da CLI para executar o Docker Compose.


---

### 📦 Provisionamento automatizado do Grafana

O Grafana utiliza o mecanismo nativo de **provisioning** através de arquivos YAML.

Dessa forma:

```text
Prometheus Datasource
        │
        ▼
Provisioning
        │
        ▼
Grafana
        │
        ▼
Dashboard automaticamente disponível
```

Isso garante que o ambiente possa ser recriado de forma consistente e sem configurações manuais.

---

## 🧪 Validação do Ambiente

Após executar o playbook, recomenda-se validar os componentes:

### Aplicação

```bash
curl http://localhost/projeto-korp
```

### Containers

```bash
docker ps
```

### Logs

```bash
docker compose logs
```

Ou, a partir do diretório do projeto:

```bash
docker compose -f docker/docker-compose.yml logs
```

---

## 🛑 Parando o ambiente

Para interromper os containers:

```bash
docker compose -f docker/docker-compose.yml down
```

Para interromper e remover também os volumes:

```bash
docker compose -f docker/docker-compose.yml down -v
```

> **Atenção:** remover os volumes pode apagar dados persistentes dos serviços que utilizam armazenamento Docker.

---

## 👨‍💻 Autor

Desenvolvido por **Breno Pombo** como solução para o desafio DevOps da **Korp**.

- **LinkedIn:** https://www.linkedin.com/in/brenopombo/
- **GitHub:** https://github.com/PBreno

---

## 🚀 Tecnologias utilizadas

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![NGINX](https://img.shields.io/badge/NGINX-009639?style=for-the-badge&logo=nginx&logoColor=white)
![Ansible](https://img.shields.io/badge/Ansible-EE0000?style=for-the-badge&logo=ansible&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white)

---

⭐ **Se este projeto foi útil ou interessante, considere deixar uma estrela no repositório.**