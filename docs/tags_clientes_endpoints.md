# Endpoints de Tags de Clientes

Este documento descreve os endpoints para gerenciar tags de clientes na API.

## Endpoints Disponíveis

### 1. Atribuir Tags a um Cliente
**POST** `/api/clientes/:id/tags`

Atribui uma ou mais tags a um cliente específico.

#### Parâmetros da URL
- `id` (int): ID do cliente

#### Body da Requisição
```json
{
  "ids_tags": [1, 2, 3]
}
```

#### Campos Obrigatórios
- `ids_tags`: Array de IDs das tags (mínimo 1 item)

#### Resposta de Sucesso (200)
```json
{
  "message": "Tags atribuídas ao cliente com sucesso"
}
```

#### Exemplo de Uso
```bash
curl -X POST "http://localhost:8080/api/clientes/1/tags" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "ids_tags": [1, 2, 3]
  }'
```

---

### 2. Remover Tags de um Cliente
**DELETE** `/api/clientes/:id/tags`

Remove uma ou mais tags específicas de um cliente.

#### Parâmetros da URL
- `id` (int): ID do cliente

#### Body da Requisição
```json
{
  "ids_tags": [1, 2]
}
```

#### Campos Obrigatórios
- `ids_tags`: Array de IDs das tags a serem removidas (mínimo 1 item)

#### Resposta de Sucesso (200)
```json
{
  "message": "Tags removidas do cliente com sucesso"
}
```

#### Exemplo de Uso
```bash
curl -X DELETE "http://localhost:8080/api/clientes/1/tags" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "ids_tags": [1, 2]
  }'
```

---

### 3. Listar Tags de um Cliente
**GET** `/api/clientes/:id/tags`

Retorna todas as tags atribuídas a um cliente específico.

#### Parâmetros da URL
- `id` (int): ID do cliente

#### Resposta de Sucesso (200)
```json
{
  "tags": [
    {
      "id": 1,
      "id_tag": 1,
      "nome": "VIP"
    },
    {
      "id": 2,
      "id_tag": 2,
      "nome": "Fidelidade"
    }
  ],
  "cliente_id": 1,
  "total": 2
}
```

#### Exemplo de Uso
```bash
curl -X GET "http://localhost:8080/api/clientes/1/tags" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Códigos de Erro

### 400 - Bad Request
- Dados de entrada inválidos
- ID do cliente inválido
- Array de tags vazio

### 401 - Unauthorized
- Token de autenticação ausente ou inválido

### 404 - Not Found
- Cliente não encontrado

### 500 - Internal Server Error
- Erro interno do servidor

---

## Estrutura da Tabela

A tabela `tags_clientes` relaciona clientes com suas tags:

```sql
CREATE TABLE `tags_clientes` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `id_tag` INT(11) NOT NULL,
  `cliente_id` INT(11) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id_tag`) REFERENCES `tags` (`id_tag`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`cliente_id`) REFERENCES `clientes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

---

## Notas Importantes

1. **Duplicatas**: O sistema evita automaticamente a criação de associações duplicadas entre cliente e tag.

2. **Transações**: As operações são realizadas de forma segura, garantindo a integridade dos dados.

3. **Validação**: Todos os endpoints validam se o cliente existe antes de realizar as operações.

4. **Autenticação**: Todos os endpoints requerem autenticação JWT válida.

5. **Logs**: Todas as operações são registradas nos logs do sistema para auditoria.