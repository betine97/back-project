# 🧪 Testes dos Endpoints de Produtos - Postman

## 🔐 Pré-requisitos

1. **Servidor rodando**: `cd cmd/server && go run main.go`
2. **Token JWT**: Obtenha fazendo login primeiro

### Obter Token JWT:
```json
POST http://localhost:8080/login
Content-Type: application/json

{
  "email": "joao.silva@email.com",
  "password": "MinhaSenh@123"
}
```

**Resposta:**
```json
{
  "message": "Login successfully",
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 📦 Testes dos Endpoints

### ✅ Teste 1: Listar Todos os Produtos

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

**Resposta Esperada (200):**
```json
{
  "products": [
    {
      "id": 1,
      "codigo_barra": "7891234560001",
      "nome_produto": "Ração Premium Cães Adultos",
      "sku": "RAC-PREM-CAES-ADU",
      "categoria": "Alimentação",
      "destinado_para": "Cães",
      "variacao": "Raças Grandes",
      "marca": "Golden",
      "descricao": "Ração completa para cães adultos de grande porte.",
      "status": "ativo",
      "preco_venda": 189.9
    },
    {
      "id": 2,
      "codigo_barra": "7891234560002",
      "nome_produto": "Ração Premium Gatos Filhotes",
      "sku": "RAC-PREM-GAT-FIL",
      "categoria": "Alimentação",
      "destinado_para": "Gatos",
      "variacao": "Filhotes",
      "marca": "Whiskas",
      "descricao": "Ração seca especialmente formulada para gatos filhotes.",
      "status": "ativo",
      "preco_venda": 99.9
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 10
}
```

### ✅ Teste 2: Buscar Produto por ID

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/1`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

**Resposta Esperada (200):**
```json
{
  "id": 1,
  "codigo_barra": "7891234560001",
  "nome_produto": "Ração Premium Cães Adultos",
  "sku": "RAC-PREM-CAES-ADU",
  "categoria": "Alimentação",
  "destinado_para": "Cães",
  "variacao": "Raças Grandes",
  "marca": "Golden",
  "descricao": "Ração completa para cães adultos de grande porte.",
  "status": "ativo",
  "preco_venda": 189.9
}
```

### ✅ Teste 3: Buscar Produto Inexistente

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/999`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

**Resposta Esperada (404):**
```json
{
  "error": "Product not found"
}
```

### ✅ Teste 4: Filtrar por Categoria

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?categoria=Alimentação`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

**Resposta Esperada (200):**
```json
{
  "products": [
    {
      "id": 1,
      "nome_produto": "Ração Premium Cães Adultos",
      "categoria": "Alimentação",
      "preco_venda": 189.9
    },
    {
      "id": 2,
      "nome_produto": "Ração Premium Gatos Filhotes",
      "categoria": "Alimentação",
      "preco_venda": 99.9
    }
  ],
  "total": 3,
  "page": 1,
  "limit": 10
}
```

### ✅ Teste 5: Filtrar por Destinação

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?destinado_para=Cães`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

### ✅ Teste 6: Filtrar por Faixa de Preço

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?min_price=50&max_price=100`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

### ✅ Teste 7: Busca Textual

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?search=ração`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

### ✅ Teste 8: Filtros Múltiplos com Paginação

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?categoria=Alimentação&destinado_para=Cães&page=1&limit=2`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

### ✅ Teste 9: Filtrar por Marca

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?marca=Golden`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

### ✅ Teste 10: Filtrar por Status

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/search?status=ativo`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

## ❌ Testes de Erro

### 🚫 Teste 11: Sem Token de Autenticação

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos`  
**Headers:**
```
Content-Type: application/json
```

**Resposta Esperada (401):**
```json
{
  "error": "Missing or malformed JWT"
}
```

### 🚫 Teste 12: Token Inválido

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos`  
**Headers:**
```
Authorization: Bearer token_invalido
Content-Type: application/json
```

**Resposta Esperada (401):**
```json
{
  "error": "Invalid or expired JWT"
}
```

### 🚫 Teste 13: ID Inválido

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos/abc`  
**Headers:**
```
Authorization: Bearer {seu_token_aqui}
Content-Type: application/json
```

**Resposta Esperada (400):**
```json
{
  "error": "Invalid product ID"
}
```

## 🎯 Sequência Recomendada de Testes

1. **Faça login** para obter o token JWT
2. **Teste listagem completa** (Teste 1)
3. **Teste busca por ID** (Teste 2)
4. **Teste filtros individuais** (Testes 4-7, 9-10)
5. **Teste filtros combinados** (Teste 8)
6. **Teste casos de erro** (Testes 11-13)

## 📊 Dados de Teste Disponíveis

Os seguintes produtos estão disponíveis para teste:

| ID | Nome | Categoria | Destinado Para | Marca | Preço |
|----|------|-----------|----------------|-------|-------|
| 1 | Ração Premium Cães Adultos | Alimentação | Cães | Golden | R$ 189,90 |
| 2 | Ração Premium Gatos Filhotes | Alimentação | Gatos | Whiskas | R$ 99,90 |
| 3 | Suplemento Vitamínico Cães Idosos | Saúde | Cães | Organnact | R$ 54,90 |
| 4 | Brinquedo Bola Mordedor | Lazer | Cães | Pet Games | R$ 24,90 |
| 5 | Ração para Pássaros Canário | Alimentação | Pássaros | Megazoo | R$ 29,90 |
| 6 | Antipulgas Gatos | Saúde | Gatos | Bayer | R$ 89,90 |
| 7 | Arranhador para Gatos | Lazer | Gatos | Chalesco | R$ 79,90 |
| 8 | Ração Cães Filhotes Raças Pequenas | Alimentação | Cães | Premier | R$ 159,90 |
| 9 | Suplemento para Pássaros | Saúde | Pássaros | Avitrin | R$ 19,90 |
| 10 | Cordão com Guizo para Gatos | Lazer | Gatos | Pet Flex | R$ 14,90 |

## 🔧 Dicas para Teste

1. **Salve o token**: Após o login, salve o token em uma variável do Postman
2. **Use Collections**: Organize os testes em uma collection
3. **Variáveis de ambiente**: Configure base_url como variável
4. **Teste paginação**: Experimente diferentes valores de page e limit
5. **Combine filtros**: Teste múltiplos filtros simultaneamente

## 🐛 Troubleshooting

- **401 Unauthorized**: Verifique se o token está correto e não expirou
- **404 Not Found**: Verifique se o ID do produto existe
- **400 Bad Request**: Verifique se os parâmetros estão corretos
- **500 Internal Server Error**: Verifique os logs do servidor