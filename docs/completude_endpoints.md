# API de Completude de Cadastro

## Visão Geral
Esta API permite verificar o percentual de completude do cadastro de clientes e seus pets associados.

## Endpoint

### GET /api/completude/clientes

Retorna a completude do cadastro de clientes com seus pets associados.

#### Parâmetros de Query
- `page` (opcional): Número da página (padrão: 1)
- `limit` (opcional): Número de itens por página (padrão: 10, máximo: 100)

#### Headers Obrigatórios
```
Authorization: Bearer <token_jwt>
```

#### Resposta de Sucesso (200)
```json
{
  "clientes": [
    {
      "cliente_id": 1,
      "nome_cliente": "João Silva",
      "percentual_completo": 85.0
    },
    {
      "cliente_id": 2,
      "nome_cliente": "Maria Santos",
      "percentual_completo": 92.5
    },
    {
      "cliente_id": 3,
      "nome_cliente": "Pedro Costa",
      "percentual_completo": 67.8
    }
  ],
  "total": 3,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

#### Campos Avaliados

**Cliente:**
- `tipo_cliente`: Tipo do cliente (obrigatório)
- `nome_cliente`: Nome completo (obrigatório)
- `numero_celular`: Número de celular (obrigatório)
- `sexo`: Sexo (M/F) (obrigatório)
- `email`: Email (opcional mas recomendado)
- `data_nascimento`: Data de nascimento (opcional)
- `data_cadastro`: Data do cadastro (obrigatório)

**Pet:**
- `nome_pet`: Nome do pet (obrigatório)
- `especie`: Espécie do pet (obrigatório)
- `raca`: Raça do pet (opcional)
- `porte`: Porte (Pequeno/Médio/Grande) (obrigatório)
- `data_aniversario`: Data de aniversário (opcional)
- `idade`: Idade em anos (opcional)

**Endereço:**
- `cep`: CEP (obrigatório)
- `logradouro`: Logradouro/Rua (obrigatório)
- `numero`: Número (obrigatório)
- `bairro`: Bairro (obrigatório)
- `cidade`: Cidade (obrigatório)
- `estado`: Estado (obrigatório)
- *Nota: Complemento e Observações são opcionais e não são considerados no cálculo*

#### Estrutura da Resposta Simples

A API retorna apenas as informações essenciais:

**Campos da Resposta:**
- `cliente_id`: ID único do cliente
- `nome_cliente`: Nome do cliente para identificação
- `percentual_completo`: Percentual geral de completude do cadastro

#### Cálculo do Percentual
O percentual é calculado considerando **TODOS** os campos de **TODAS** as tabelas relacionadas:

1. **Tabela clientes**: 7 campos obrigatórios
2. **Tabela pets**: 6 campos por pet (multiplicado pelo número de pets)
3. **Tabela endereços**: 6 campos por endereço (multiplicado pelo número de endereços)

**Fórmula:** `(total_campos_preenchidos / total_campos_obrigatórios) * 100`

#### Exemplo Prático
```
Cliente ID 1: João Silva
- Clientes: 7/7 campos preenchidos
- Pets: 5/6 campos preenchidos (1 pet)
- Endereços: 4/6 campos preenchidos (1 endereço)
- Total: 16/19 campos = 84.21%
```

#### Interpretação Rápida
- **90-100%**: Cadastro quase completo
- **70-89%**: Cadastro bom, poucos campos faltando
- **50-69%**: Cadastro médio, vários campos faltando
- **0-49%**: Cadastro incompleto, muitos campos faltando

#### Exemplos de Uso

**Buscar primeira página:**
```
GET /api/completude/clientes?page=1&limit=10
```

**Buscar com paginação específica:**
```
GET /api/completude/clientes?page=2&limit=5
```

#### Códigos de Erro
- `400`: Parâmetros inválidos
- `401`: Token de autenticação inválido ou ausente
- `500`: Erro interno do servidor

#### Exemplo de Resposta de Erro
```json
{
  "error": "Token de autenticação inválido"
}
```

## Casos de Uso
- Dashboard administrativo mostrando clientes com cadastros incompletos
- Relatórios de qualidade de dados
- Campanhas de atualização de cadastro
- Métricas de completude para análise de negócio