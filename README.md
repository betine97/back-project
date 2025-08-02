# Back Project

## Configuração do Ambiente

### Arquivo de Configuração

O projeto agora utiliza um único arquivo `.env` na raiz do projeto para todas as configurações. 

1. Copie o arquivo `.env.example` para `.env`:
   ```bash
   cp .env.example .env
   ```

2. Configure as variáveis de ambiente no arquivo `.env` conforme necessário.

### Executando a Aplicação

#### Opção 1: Script de Execução (Recomendado)
```bash
# Windows
run.bat

# Linux/Mac
./run.sh
```

#### Opção 2: Execução Manual
```bash
cd cmd/server
go run main.go
```

#### Opção 3: Build e Execução
```bash
# Build
cd cmd/server
go build -o ../../server.exe main.go

# Executar
cd ../..
./server.exe
```

## Estrutura do Projeto

```
.
├── .env                    # Configurações do ambiente (único arquivo)
├── .env.example           # Exemplo de configurações
├── cmd/
│   ├── config/           # Configurações da aplicação
│   └── server/           # Ponto de entrada da aplicação
├── src/
│   ├── controller/       # Controladores HTTP
│   ├── model/           # Modelos, serviços e persistência
│   └── view/            # Views (se aplicável)
├── run.bat              # Script de execução Windows
├── run.sh               # Script de execução Linux/Mac
└── server.exe           # Executável compilado
```

## Mudanças Realizadas

- ✅ Consolidado todos os arquivos `.env` em um único arquivo na raiz
- ✅ Corrigido carregamento de configurações para usar o `.env` da raiz
- ✅ Corrigido porta hardcoded no main.go para usar variável de ambiente
- ✅ Adicionado `.env` ao `.gitignore` para segurança
- ✅ Criado `.env.example` para documentação
- ✅ Criados scripts de execução para facilitar o desenvolvimento