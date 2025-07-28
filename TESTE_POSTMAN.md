# Guia de Testes - API de Usuários

## Configuração Inicial

1. Certifique-se de que o servidor está rodando na porta 8080
2. Execute o comando para criar a tabela de usuários:
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       first_name VARCHAR(100) NOT NULL,     
       last_name VARCHAR(100) NOT NULL,     
       email VARCHAR(255) NOT NULL UNIQUE,  
       city VARCHAR(100) NOT NULL,          
       password VARCHAR(255) NOT NULL        
   );
   ```

## Testes de Registro (/register)

### 1. Registro de Usuário Válido
**Método:** POST  
**URL:** `http://localhost:8080/register`  
**Headers:** `Content-Type: application/json`

**Body:**
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

### 2. Mais Exemplos de Registro

**Maria Santos:**
```json
{
  "first_name": "Maria",
  "last_name": "Santos",
  "email": "maria.santos@email.com",
  "city": "Rio de Janeiro",
  "password": "Senha#456"
}
```

**Pedro Costa:**
```json
{
  "first_name": "Pedro",
  "last_name": "Costa",
  "email": "pedro.costa@email.com",
  "city": "Belo Horizonte",
  "password": "Test@789"
}
```

## Testes de Login (/login)

### 1. Login Válido
**Método:** POST  
**URL:** `http://localhost:8080/login`  
**Headers:** `Content-Type: application/json`

**Body:**
```json
{
  "email": "joao.silva@email.com",
  "password": "MinhaSenh@123"
}
```

**Resposta Esperada (200):**
```json
{
  "message": "Login successfully",
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Testes de Validação (Casos de Erro)

### 1. Nome muito curto
```json
{
  "first_name": "A",
  "last_name": "Silva",
  "email": "teste@email.com",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```
**Erro esperado:** Validation error

### 2. Email inválido
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "email-invalido",
  "city": "São Paulo",
  "password": "MinhaSenh@123"
}
```
**Erro esperado:** Invalid email format

### 3. Senha sem caracteres especiais
```json
{
  "first_name": "João",
  "last_name": "Silva",
  "email": "joao@email.com",
  "city": "São Paulo",
  "password": "senha123"
}
```
**Erro esperado:** Password validation error

### 4. Email duplicado
Tente registrar o mesmo email duas vezes:
**Erro esperado:** "Email is already associated with an existing account"

## Testando Endpoints Protegidos

Após fazer login e obter o token, você pode testar endpoints protegidos:

**Método:** GET  
**URL:** `http://localhost:8080/api/produtos`  
**Headers:** 
- `Content-Type: application/json`
- `Authorization: Bearer SEU_TOKEN_AQUI`

## Códigos de Status Esperados

- **201:** Usuário criado com sucesso
- **200:** Login realizado com sucesso
- **400:** Erro de validação ou dados inválidos
- **401:** Credenciais inválidas
- **409:** Email já existe
- **500:** Erro interno do servidor

## Dicas para Teste

1. Teste primeiro o registro de usuários válidos
2. Em seguida, teste o login com as credenciais criadas
3. Teste casos de erro para validar as validações
4. Use o token retornado no login para acessar endpoints protegidos
5. Verifique se os dados estão sendo salvos corretamente na tabela `users`