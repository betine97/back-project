# Endpoints de Tags

Este documento descreve os endpoints para gerenciar tags na API.

## Endpoints Disponíveis

### 1. Listar Todas as Tags
**GET** `/api/tags`

Retorna todas as tags disponíveis no sistema com paginação.

#### Parâmetros de Query (Opcionais)
- `page` (int): Número da página (padrão: 1)
- `limit` (int): Itens por página (padrão: 30, máximo: 100)

#### Resposta de Sucesso (200)
```json
{
  "tags": [
    {
      "id_tag": 1,
      "categoria_tag": "Fidelidade",
      "nome_tag": "VIP"
    },
    {
      "id_tag": 2,
      "categoria_tag": "Fidelidade",
      "nome_tag": "Premium"
    },
    {
      "id_tag": 3,
      "categoria_tag": "Comportamento",
      "nome_tag": "Frequente"
    }
  ],
  "total": 3,
  "page": 1,
  "limit": 30,
  "total_pages": 1
}
```

#### Exemplo de Uso
```bash
# Buscar todas as tags (primeira página)
curl -X GET "http://localhost:8080/api/tags" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Buscar com paginação
curl -X GET "http://localhost:8080/api/tags?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 2. Criar Nova Tag
**POST** `/api/tags`

Cria uma nova tag no sistema.

#### Body da Requisição
```json
{
  "categoria_tag": "Fidelidade",
  "nome_tag": "Gold"
}
```

#### Campos Obrigatórios
- `categoria_tag`: Categoria da tag (2-100 caracteres)
- `nome_tag`: Nome da tag (2-100 caracteres)

#### Resposta de Sucesso (201)
```json
{
  "message": "Tag created successfully"
}
```

#### Exemplo de Uso
```bash
curl -X POST "http://localhost:8080/api/tags" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "categoria_tag": "Fidelidade",
    "nome_tag": "Gold"
  }'
```

---

## Exemplos de Categorias e Tags

### Sugestões de Categorias:
- **Fidelidade**: VIP, Premium, Gold, Silver, Bronze
- **Comportamento**: Frequente, Esporádico, Novo Cliente
- **Preferências**: Pet Lover, Eco-Friendly, Tech Savvy
- **Localização**: Local, Regional, Nacional
- **Valor**: Alto Valor, Médio Valor, Baixo Valor

---

## Códigos de Erro

### 400 - Bad Request
```json
{
  "request invalid": {
    "message": "Some fields are invalid",
    "causes": [
      {
        "field": "categoria_tag",
        "message": "categoria_tag is a required field"
      }
    ]
  }
}
```

### 401 - Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 500 - Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Integração com Tags de Clientes

Após criar tags usando estes endpoints, você pode usar os IDs retornados para:

1. **Atribuir tags aos clientes**: `POST /api/clientes/:id/tags`
2. **Remover tags dos clientes**: `DELETE /api/clientes/:id/tags`
3. **Listar tags dos clientes**: `GET /api/clientes/:id/tags`

### Exemplo de Fluxo Completo:

```bash
# 1. Criar uma nova tag
curl -X POST "http://localhost:8080/api/tags" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"categoria_tag": "Fidelidade", "nome_tag": "Platinum"}'

# 2. Listar todas as tags para ver o ID da nova tag
curl -X GET "http://localhost:8080/api/tags" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. Atribuir a nova tag ao cliente (assumindo que o ID da tag é 5)
curl -X POST "http://localhost:8080/api/clientes/1/tags" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ids_tags": [5]}'
```

---

## Estrutura da Tabela

```sql
CREATE TABLE `tags` (
  `id_tag` int(11) NOT NULL AUTO_INCREMENT,
  `categoria_tag` varchar(100) NOT NULL,
  `nome_tag` varchar(100) NOT NULL,
  PRIMARY KEY (`id_tag`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

---

## Notas Importantes

1. **Autenticação**: Todos os endpoints requerem autenticação JWT válida.

2. **Paginação**: O endpoint GET suporta paginação para melhor performance.

3. **Validação**: Todos os campos são validados antes da criação.

4. **Logs**: Todas as operações são registradas nos logs do sistema.

5. **Ordenação**: As tags são retornadas ordenadas por ID decrescente (mais recentes primeiro).