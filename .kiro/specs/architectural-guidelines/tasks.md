# Implementation Plan

Este documento serve como guia de referência para implementação de novas funcionalidades seguindo os padrões arquiteturais estabelecidos.

## Checklist para Implementação de Nova Funcionalidade

- [ ] 1. Análise e Planejamento
  - Verificar se a funcionalidade se alinha com a arquitetura MVC existente
  - Identificar quais camadas serão afetadas (Controller/Service/Model)
  - Definir interfaces necessárias para desacoplamento
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [ ] 2. Implementação da Camada de Modelo
- [ ] 2.1 Criar entidades de domínio
  - Implementar structs em `src/model/entity/` seguindo padrão existente
  - Incluir tags GORM e JSON apropriadas
  - Implementar método `NewID()` se necessário
  - _Requirements: 1.3, 4.3_

- [ ] 2.2 Criar DTOs de entrada e saída
  - Implementar structs em `src/model/dtos/` com validações
  - Utilizar tags `validate` do go-playground/validator
  - Separar DTOs de request e response quando necessário
  - _Requirements: 2.4, 4.4, 6.2_

- [ ] 2.3 Implementar camada de persistência
  - Criar interface em `src/model/persistence/`
  - Implementar métodos CRUD utilizando GORM
  - Seguir padrão de tratamento de erros estabelecido
  - _Requirements: 2.2, 3.3, 4.5_

- [ ] 3. Implementação da Camada de Serviço
- [ ] 3.1 Definir interface do serviço
  - Criar interface específica para o domínio
  - Definir métodos que encapsulam lógica de negócio
  - Considerar dependências necessárias (crypto, persistence)
  - _Requirements: 1.4, 2.1_

- [ ] 3.2 Implementar lógica de negócio
  - Criar struct que implementa a interface do serviço
  - Implementar validações de negócio
  - Utilizar crypto service para operações sensíveis
  - Implementar transformações entre DTOs e entidades
  - _Requirements: 2.1, 2.5, 6.1_

- [ ] 3.3 Adicionar tratamento de erros
  - Utilizar padrão RestErr para erros estruturados
  - Implementar validações com retorno de erros apropriados
  - Adicionar logs estruturados com Zap
  - _Requirements: 3.1, 3.2, 3.4_

- [ ] 4. Implementação da Camada de Controller
- [ ] 4.1 Criar métodos do controller
  - Implementar handlers em `src/controller/`
  - Seguir padrão de recebimento de *fiber.Ctx
  - Implementar binding e validação de entrada
  - _Requirements: 1.1, 4.1, 6.1_

- [ ] 4.2 Integrar com camada de serviço
  - Utilizar injeção de dependências estabelecida
  - Chamar métodos do serviço apropriados
  - Transformar respostas do serviço em JSON
  - _Requirements: 1.2, 2.1_

- [ ] 4.3 Implementar tratamento de respostas
  - Retornar códigos HTTP apropriados
  - Utilizar fiber.Map para respostas JSON
  - Tratar erros do serviço adequadamente
  - _Requirements: 3.1, 3.2_

- [ ] 5. Configuração de Rotas
- [ ] 5.1 Adicionar rotas em routes package
  - Registrar novos endpoints seguindo padrão existente
  - Configurar middlewares se necessário
  - Documentar endpoints adequadamente
  - _Requirements: 1.1, 6.3_

- [ ] 6. Configuração e Integração
- [ ] 6.1 Atualizar injeção de dependências
  - Modificar `initDependencies()` em main.go se necessário
  - Manter padrão de inicialização estabelecido
  - Verificar configurações de ambiente necessárias
  - _Requirements: 1.2, 5.3_

- [ ] 6.2 Adicionar configurações de ambiente
  - Incluir novas variáveis no .env se necessário
  - Seguir padrão de nomenclatura estabelecido
  - Documentar configurações no README
  - _Requirements: 5.1, 5.2_

- [ ] 7. Testes e Validação
- [ ] 7.1 Implementar testes unitários
  - Criar testes para cada método do serviço
  - Implementar mocks para dependências
  - Seguir padrão Arrange-Act-Assert
  - _Requirements: 6.4_

- [ ] 7.2 Implementar testes de integração
  - Testar endpoints completos
  - Validar fluxos de dados entre camadas
  - Testar cenários de erro
  - _Requirements: 6.4_

- [ ] 7.3 Validar implementação
  - Verificar se segue todos os padrões estabelecidos
  - Testar com diferentes cenários de entrada
  - Validar logs e tratamento de erros
  - _Requirements: 6.5_

## Padrões de Referência Rápida

### Estrutura de Controller
```go
func (ctl *Controller) NewMethod(ctx *fiber.Ctx) error {
    var request dtos.RequestDTO
    if err := ctx.BodyParser(&request); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(exceptions.NewBadRequestError("Invalid request body"))
    }
    
    result, err := ctl.service.ServiceMethod(request)
    if err != nil {
        return ctx.Status(err.Code).JSON(err)
    }
    
    return ctx.Status(fiber.StatusOK).JSON(result)
}
```

### Estrutura de Service
```go
func (s *Service) ServiceMethod(request dtos.RequestDTO) (*dtos.ResponseDTO, *exceptions.RestErr) {
    // Validações de negócio
    if err := s.validateRequest(request); err != nil {
        return nil, err
    }
    
    // Transformação para entidade
    entity := buildEntity(request)
    
    // Operação de persistência
    result, err := s.persistence.Operation(entity)
    if err != nil {
        return nil, exceptions.NewInternalServerError("Database operation failed")
    }
    
    // Transformação para DTO de resposta
    return buildResponseDTO(result), nil
}
```

### Estrutura de Persistence
```go
func (p *Persistence) Operation(entity *entity.Entity) (*entity.Entity, error) {
    if err := p.db.Create(entity).Error; err != nil {
        return nil, err
    }
    return entity, nil
}
```