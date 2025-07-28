# 📦 Endpoints de Produtos - API Documentation

## 🎯 Endpoints Implementados

### 1. **GET /api/produtos** - Listar Todos os Produtos
Retorna todos os produtos cadastrados no sistema.

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Resposta (200):**
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
      "preco_venda": 189.90
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 10
}
```

### 2. **GET /api/produtos/:id** - Buscar Produto por ID
Retorna um produto específico pelo ID.

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Exemplo:**
```
GET /api/produtos/1
```

**Resposta (200):**
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
  "preco_venda": 189.90
}
```

**Resposta (404):**
```json
{
  "error": "Product not found"
}
```

### 3. **GET /api/produtos/search** - Buscar Produtos com Filtros
Permite buscar produtos com diversos filtros e paginação.

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Query Parameters:**
- `page` (int, opcional): Número da página (padrão: 1)
- `limit` (int, opcional): Itens por página (padrão: 10)
- `categoria` (string, opcional): Filtrar por categoria
- `destinado_para` (string, opcional): Filtrar por destinação
- `marca` (string, opcional): Filtrar por marca
- `variacao` (string, opcional): Filtrar por variação
- `status` (string, opcional): Filtrar por status (padrão: "ativo")
- `min_price` (float, opcional): Preço mínimo
- `max_price` (float, opcional): Preço máximo
- `search` (string, opcional): Busca textual em nome, descrição e SKU

**Exemplos de Uso:**

#### Buscar produtos para cães:
```
GET /api/produtos/search?destinado_para=Cães
```

#### Buscar produtos da categoria Alimentação:
```
GET /api/produtos/search?categoria=Alimentação
```

#### Buscar produtos com preço entre R$ 50 e R$ 100:
```
GET /api/produtos/search?min_price=50&max_price=100
```

#### Buscar produtos com texto "ração":
```
GET /api/produtos/search?search=ração
```

#### Buscar com múltiplos filtros e paginação:
```
GET /api/produtos/search?categoria=Alimentação&destinado_para=Cães&page=1&limit=5
```

**Resposta (200):**
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
      "preco_venda": 189.90
    }
  ],
  "total": 3,
  "page": 1,
  "limit": 10
}
```

## 🧪 Exemplos de Teste no Postman

### 1. Teste Básico - Listar Todos os Produtos
```
Method: GET
URL: http://localhost:8080/api/produtos
Headers:
  Authorization: Bearer {seu_token_jwt}
  Content-Type: application/json
```

### 2. Teste - Buscar Produto por ID
```
Method: GET
URL: http://localhost:8080/api/produtos/1
Headers:
  Authorization: Bearer {seu_token_jwt}
  Content-Type: application/json
```

### 3. Teste - Filtros Avançados
```
Method: GET
URL: http://localhost:8080/api/produtos/search?categoria=Alimentação&destinado_para=Cães&min_price=50&max_price=200&page=1&limit=5
Headers:
  Authorization: Bearer {seu_token_jwt}
  Content-Type: application/json
```

## 🔐 Autenticação

Todos os endpoints de produtos requerem autenticação JWT. Para obter o token:

1. **Registre um usuário:**
```json
POST /register
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "joao@email.com",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```

2. **Faça login:**
```json
POST /login
{
  "email": "joao@email.com",
  "password": "MinhaSenh@123"
}
```

3. **Use o token retornado** nos headers dos endpoints de produtos.

## 📊 Códigos de Status

| Código | Descrição |
|--------|-----------|
| 200 | Sucesso |
| 400 | Parâmetros inválidos |
| 401 | Token JWT inválido ou ausente |
| 404 | Produto não encontrado |
| 500 | Erro interno do servidor |

## 🎨 Integração com Frontend

### JavaScript/Fetch Example:
```javascript
// Buscar todos os produtos
const response = await fetch('http://localhost:8080/api/produtos', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();

// Buscar produtos com filtros
const filteredResponse = await fetch(
  'http://localhost:8080/api/produtos/search?categoria=Alimentação&destinado_para=Cães',
  {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  }
);
const filteredData = await filteredResponse.json();
```

### React Example:
```jsx
const [products, setProducts] = useState([]);
const [loading, setLoading] = useState(true);

useEffect(() => {
  const fetchProducts = async () => {
    try {
      const response = await fetch('/api/produtos', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      const data = await response.json();
      setProducts(data.products);
    } catch (error) {
      console.error('Error fetching products:', error);
    } finally {
      setLoading(false);
    }
  };

  fetchProducts();
}, [token]);
```

## ✨ Recursos Implementados

- ✅ **Listagem completa** de produtos
- ✅ **Busca por ID** específico
- ✅ **Filtros avançados** (categoria, marca, preço, etc.)
- ✅ **Busca textual** em nome, descrição e SKU
- ✅ **Paginação** configurável
- ✅ **Autenticação JWT** obrigatória
- ✅ **Logs estruturados** para monitoramento
- ✅ **Tratamento de erros** adequado
- ✅ **Validação de parâmetros**
- ✅ **Resposta padronizada** com metadados

## 🚀 Próximos Passos

1. **Teste os endpoints** usando o Postman
2. **Integre com seu frontend** React/Vue/Angular
3. **Implemente cache** para melhor performance
4. **Adicione mais filtros** conforme necessário