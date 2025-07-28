# 🚀 Massa de Teste - API de Usuários

## ⚙️ Configuração Inicial

1. **Inicie o servidor:**
   ```bash
   ./server.exe
   ```

2. **Verifique se o servidor está rodando:**
   - URL base: `http://localhost:8080`
   - Porta: `8080`

## 📝 Testes de Registro de Usuários

### ✅ Teste 1: Registro Válido - João Silva
**Método:** `POST`  
**URL:** `http://localhost:8080/register`  
**Headers:** `Content-Type: application/json`

```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "joao.silva@email.com",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```

**Resposta Esperada (201):**
```json
{
  "message": "User created successfully",
  "usuário": {
    "id": 1,
    "first_name": "João",
    "last_name": "Silva",
    "email": "joao.silva@email.com",
    "city": "São Paulo"
  }
}
```

### ✅ Teste 2: Registro Válido - Maria Santos
```json
{
  "first_name": "Maria",
  "last_name": "Santos",
  "email": "maria.santos@email.com",
  "city": "Rio de Janeiro",
  "password": "Senha#456"
}
```

### ✅ Teste 3: Registro Válido - Pedro Costa
```json
{
  "first_name": "Pedro",
  "last_name": "Costa",
  "email": "pedro.costa@email.com",
  "city": "Belo Horizonte",
  "password": "Test@789"
}
```

### ✅ Teste 4: Registro Válido - Ana Oliveira
```json
{
  "first_name": "Ana",
  "last_name": "Oliveira",
  "email": "ana.oliveira@email.com",
  "city": "Brasília",
  "password": "Ana@2024"
}
```

### ✅ Teste 5: Registro Válido - Carlos Ferreira
```json
{
  "first_name": "Carlos",
  "last_name": "Ferreira",
  "email": "carlos.ferreira@email.com",
  "city": "Salvador",
  "password": "Carlos#789"
}
```

## 🔐 Testes de Login

### ✅ Login 1: João Silva
**Método:** `POST`  
**URL:** `http://localhost:8080/login`  
**Headers:** `Content-Type: application/json`

```json
{
  "email": "joao.silva@email.com",
  "password": "MinhaSenh@123"
}
```

### ✅ Login 2: Maria Santos
```json
{
  "email": "maria.santos@email.com",
  "password": "Senha#456"
}
```

### ✅ Login 3: Pedro Costa
```json
{
  "email": "pedro.costa@email.com",
  "password": "Test@789"
}
```

**Resposta Esperada do Login (200):**
```json
{
  "message": "Login successfully",
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## ❌ Testes de Validação (Casos de Erro)

### 🚫 Erro 1: Nome muito curto
```json
{
  "first_name": "A",
  "last_name": "Silva",
  "email": "teste1@email.com",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```

### 🚫 Erro 2: Sobrenome muito curto
```json
{
  "first_name": "João",
  "last_name": "S",
  "email": "teste2@email.com",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```

### 🚫 Erro 3: Email inválido
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "email-sem-arroba",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```

### 🚫 Erro 4: Cidade muito curta
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "teste4@email.com",
  "city": "A",
  "password": "MinhaSenh@123"
}
```

### 🚫 Erro 5: Senha sem caracteres especiais
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "teste5@email.com",
  "city": "São Paulo",
  "password": "senha123"
}
```

### 🚫 Erro 6: Senha muito curta
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "teste6@email.com",
  "city": "São Paulo",
  "password": "123@"
}
```

### 🚫 Erro 7: Email duplicado
Tente registrar novamente com um email já usado:
```json
{
  "first_name": "Outro",
  "last_name": "Usuario",
  "email": "joao.silva@email.com",
  "city": "Recife",
  "password": "OutraSenha@123"
}
```

### 🚫 Erro 8: Campos obrigatórios ausentes
```json
{
  "first_name": "João",
  "email": "teste8@email.com",
  "password": "MinhaSenh@123"
}
```

### 🚫 Erro 9: Login com email inexistente
```json
{
  "email": "naoexiste@email.com",
  "password": "qualquersenha"
}
```

### 🚫 Erro 10: Login com senha incorreta
```json
{
  "email": "joao.silva@email.com",
  "password": "senhaerrada"
}
```

## 🔒 Teste de Endpoint Protegido

Após fazer login e obter o token, teste o endpoint protegido:

**Método:** `GET`  
**URL:** `http://localhost:8080/api/produtos`  
**Headers:** 
- `Content-Type: application/json`
- `Authorization: Bearer SEU_TOKEN_AQUI`

## 📊 Códigos de Status Esperados

| Código | Descrição |
|--------|-----------|
| 201 | Usuário criado com sucesso |
| 200 | Login realizado com sucesso |
| 400 | Erro de validação ou dados inválidos |
| 401 | Credenciais inválidas |
| 409 | Email já existe |
| 500 | Erro interno do servidor |

## 🎯 Sequência Recomendada de Testes

1. **Registre 3-5 usuários válidos** (Testes 1-5)
2. **Teste login com os usuários criados** (Login 1-3)
3. **Teste casos de validação** (Erros 1-8)
4. **Teste login com credenciais inválidas** (Erros 9-10)
5. **Use o token para acessar endpoint protegido**

## 🔧 Dicas Importantes

- ✅ Sempre use `Content-Type: application/json` no header
- ✅ Senhas devem ter pelo menos 6 caracteres e conter caracteres especiais
- ✅ Nomes e sobrenomes devem ter pelo menos 2 caracteres
- ✅ Emails devem ter formato válido
- ✅ Cidades devem ter pelo menos 2 caracteres
- ✅ Guarde o token retornado no login para usar em endpoints protegidos
- ✅ Cada email só pode ser registrado uma vez

## 🐛 Troubleshooting

Se encontrar erros:
1. Verifique se o servidor está rodando na porta 8080
2. Confirme que os headers estão corretos
3. Valide o formato JSON do body
4. Verifique se não há campos extras no JSON
5. Confirme que a tabela `users` foi criada corretamente