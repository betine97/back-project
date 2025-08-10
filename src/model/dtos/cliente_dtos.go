package dtos

// Para GET api/clientes
type ClienteResponse struct {
	ID             int    `json:"id"`
	TipoCliente    string `json:"tipo_cliente"`
	NomeCliente    string `json:"nome_cliente"`
	NumeroCelular  string `json:"numero_celular"`
	Sexo           string `json:"sexo"`
	Email          string `json:"email"`
	DataNascimento string `json:"data_nascimento"`
	DataCadastro   string `json:"data_cadastro"`
}

type ClienteListResponse struct {
	Clientes   []ClienteResponse `json:"clientes"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// Para POST api/clientes
type CreateClienteRequest struct {
	TipoCliente    string `json:"tipo_cliente" validate:"required,max=11"`
	NomeCliente    string `json:"nome_cliente" validate:"required,min=2,max=100"`
	NumeroCelular  string `json:"numero_celular" validate:"required,max=15"`
	Sexo           string `json:"sexo" validate:"required,oneof=M F"`
	Email          string `json:"email" validate:"email,max=100"`
	DataNascimento string `json:"data_nascimento" validate:"max=20"`
	DataCadastro   string `json:"data_cadastro" validate:"required,max=20"`
}

// Para POST api/clientes/:id/tags - Atribuir tags a um cliente
type AtribuirTagsClienteRequest struct {
	IDsTags []int `json:"ids_tags" validate:"required,min=1"`
}

// Para DELETE api/clientes/:id/tags - Remover tags de um cliente
type RemoverTagsClienteRequest struct {
	IDsTags []int `json:"ids_tags" validate:"required,min=1"`
}

// Para GET api/clientes/:id/tags - Listar tags de um cliente
type TagClienteResponse struct {
	ID    int    `json:"id"`
	IDTag int    `json:"id_tag"`
	Nome  string `json:"nome"`
}

type TagsClienteListResponse struct {
	Tags      []TagClienteResponse `json:"tags"`
	ClienteID int                  `json:"cliente_id"`
	Total     int                  `json:"total"`
}

// Para GET api/clientes/buscar-criterios - Buscar clientes para validação de critérios
type ClienteCriterioResponse struct {
	ID          int    `json:"id"`
	TipoCliente string `json:"tipo_cliente"`
	Sexo        string `json:"sexo"`
}

type ClienteCriterioListResponse struct {
	Clientes []ClienteCriterioResponse `json:"clientes"`
	Total    int                       `json:"total"`
}

// Para POST api/clientes/adicionar-ao-publico/:id_publico - Resposta da adição de clientes ao público
type AdicionarClientesPublicoResponse struct {
	ClientesAdicionados int                       `json:"clientes_adicionados"`
	ClientesJaExistiam  int                       `json:"clientes_ja_existiam"`
	ClientesEncontrados []ClienteCriterioResponse `json:"clientes_encontrados"`
	Total               int                       `json:"total"`
}

// Para GET api/publicos/:id/clientes - Listar clientes de um público
type ClientePublicoResponse struct {
	ID             int    `json:"id"`
	TipoCliente    string `json:"tipo_cliente"`
	NomeCliente    string `json:"nome_cliente"`
	NumeroCelular  string `json:"numero_celular"`
	Sexo           string `json:"sexo"`
	Email          string `json:"email"`
	DataNascimento string `json:"data_nascimento"`
	DataCadastro   string `json:"data_cadastro"`
	DataAdicao     string `json:"data_adicao"`
}

type ClientesPublicoListResponse struct {
	Clientes   []ClientePublicoResponse `json:"clientes"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
	IDPublico  int                      `json:"id_publico"`
}
