package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/betine97/back-project.git/cmd/config/exceptions"
	dtos "github.com/betine97/back-project.git/src/model/dtos"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	en := en.New()
	unt := ut.New(en, en)
	transl, _ = unt.GetTranslator("en")
	en_translation.RegisterDefaultTranslations(Validate, transl)
}

func UserValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting user validation")

	var createUser dtos.CreateUser
	data := ctx.Body()

	err := ValidateUnexpectedFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createUser); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createUser); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createUser", createUser)
	zap.L().Info("User validation completed successfully", zap.Error(err))
	return ctx.Next()
}

func ProductValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting product validation")

	var createProduct dtos.CreateProductRequest
	data := ctx.Body()

	err := ValidateUnexpectedProductFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createProduct); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createProduct); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createProduct", createProduct)
	zap.L().Info("Product validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedFields(ctx *fiber.Ctx, data []byte) error {

	zap.L().Info("Validating unexpected fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"first_name":   true,
		"last_name":    true,
		"email":        true,
		"nome_empresa": true, // Novo campo
		"categoria":    true, // Novo campo
		"segmento":     true, // Novo campo
		"city":         true,
		"state":        true, // Adicionei o campo state também
		"password":     true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})

}

func PedidoValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting pedido validation")

	var createPedido dtos.CreatePedidoRequest
	data := ctx.Body()

	err := ValidateUnexpectedPedidoFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createPedido); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createPedido); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createPedido", createPedido)
	zap.L().Info("Pedido validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedProductFields(ctx *fiber.Ctx, data []byte) error {

	zap.L().Info("Validating unexpected product fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"data_cadastro":  true,
		"codigo_barra":   true,
		"nome_produto":   true,
		"sku":            true,
		"categoria":      true,
		"destinado_para": true,
		"variacao":       true,
		"marca":          true,
		"descricao":      true,
		"status":         true,
		"preco_venda":    true,
		"id_fornecedor":  true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected product fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})

}

func ItemPedidoValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting item pedido validation")

	var createItemPedido dtos.CreateItemPedidoRequest
	data := ctx.Body()

	err := ValidateUnexpectedItemPedidoFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createItemPedido); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createItemPedido); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createItemPedido", createItemPedido)
	zap.L().Info("Item pedido validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedPedidoFields(ctx *fiber.Ctx, data []byte) error {

	zap.L().Info("Validating unexpected pedido fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"id_fornecedor":    true,
		"data_pedido":      true,
		"data_entrega":     true,
		"valor_frete":      true,
		"custo_pedido":     true,
		"valor_total":      true,
		"descricao_pedido": true,
		"status":           true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected pedido fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})

}

func ValidateUnexpectedItemPedidoFields(ctx *fiber.Ctx, data []byte) error {

	zap.L().Info("Validating unexpected item pedido fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"id_produto":     true,
		"quantidade":     true,
		"preco_unitario": true,
		"subtotal":       true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected item pedido fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})

}
func EstoqueValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting estoque validation")

	var createEstoque dtos.CreateEstoqueRequest
	data := ctx.Body()

	err := ValidateUnexpectedEstoqueFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createEstoque); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createEstoque); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createEstoque", createEstoque)
	zap.L().Info("Estoque validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedEstoqueFields(ctx *fiber.Ctx, data []byte) error {
	zap.L().Info("Validating unexpected estoque fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"id_produto":           true,
		"id_lote":              true,
		"quantidade":           true,
		"vencimento":           true,
		"custo_unitario":       true,
		"data_entrada":         true,
		"data_saida":           true,
		"documento_referencia": true,
		"status":               true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected estoque fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})
}

func ClienteValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting cliente validation")

	var createCliente dtos.CreateClienteRequest
	data := ctx.Body()

	err := ValidateUnexpectedClienteFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createCliente); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createCliente); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createCliente", createCliente)
	zap.L().Info("Cliente validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedClienteFields(ctx *fiber.Ctx, data []byte) error {
	zap.L().Info("Validating unexpected cliente fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"tipo_cliente":    true,
		"nome_cliente":    true,
		"numero_celular":  true,
		"sexo":            true,
		"email":           true,
		"data_nascimento": true,
		"data_cadastro":   true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected cliente fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})
}
func EnderecoValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting endereco validation")

	var createEndereco dtos.CreateEnderecoRequest
	data := ctx.Body()

	err := ValidateUnexpectedEnderecoFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createEndereco); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createEndereco); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createEndereco", createEndereco)
	zap.L().Info("Endereco validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedEnderecoFields(ctx *fiber.Ctx, data []byte) error {
	zap.L().Info("Validating unexpected endereco fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"id_cliente":  true,
		"cep":         true,
		"cidade":      true,
		"estado":      true,
		"bairro":      true,
		"logradouro":  true,
		"numero":      true,
		"complemento": true,
		"obs":         true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected endereco fields")
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": fmt.Sprintf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields),
	})
}
func PublicoValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting publico validation")

	var createPublico dtos.CreatePublicoRequest
	data := ctx.Body()

	err := ValidateUnexpectedPublicoFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createPublico); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createPublico); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createPublico", createPublico)
	zap.L().Info("Publico validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedPublicoFields(ctx *fiber.Ctx, data []byte) error {
	zap.L().Info("Validating unexpected publico fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"nome":         true,
		"descricao":    true,
		"data_criacao": true,
		"status":       true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected publico fields")
	return fmt.Errorf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields)
}
func PetValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting pet validation")

	var createPet dtos.CreatePetRequest
	data := ctx.Body()

	err := ValidateUnexpectedPetFields(ctx, data)
	if err != nil {
		zap.L().Error("Unexpected fields in the request", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := json.Unmarshal(data, &createPet); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(&createPet); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createPet", createPet)
	zap.L().Info("Pet validation completed successfully")
	return ctx.Next()
}

func ValidateUnexpectedPetFields(ctx *fiber.Ctx, data []byte) error {
	zap.L().Info("Validating unexpected pet fields")

	var rawMap map[string]interface{}

	if err := json.Unmarshal(data, &rawMap); err != nil {
		zap.L().Error("Formato de JSON inválido", zap.Error(err))
		return exceptions.NewBadRequestError("Invalid JSON format")
	}

	expectedFields := map[string]bool{
		"cliente_id":       true,
		"nome_pet":         true,
		"especie":          true,
		"raca":             true,
		"porte":            true,
		"data_aniversario": true,
		"idade":            true,
	}

	var unexpectedFields []string
	for field := range rawMap {
		if !expectedFields[field] {
			unexpectedFields = append(unexpectedFields, field)
		}
	}

	if len(unexpectedFields) == 0 {
		return nil
	}

	zap.L().Info("Validating unexpected pet fields")
	return fmt.Errorf("Unexpected fields: %v. Please remove them and try again.", unexpectedFields)
}
func AtribuirTagsClienteValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting atribuir tags cliente validation")

	var request dtos.AtribuirTagsClienteRequest
	data := ctx.Body()

	if err := json.Unmarshal(data, &request); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(request); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("atribuirTagsCliente", request)
	return ctx.Next()
}

func RemoverTagsClienteValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting remover tags cliente validation")

	var request dtos.RemoverTagsClienteRequest
	data := ctx.Body()

	if err := json.Unmarshal(data, &request); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(request); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("removerTagsCliente", request)
	return ctx.Next()
}
func TagValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting tag validation")

	var createTag dtos.CreateTagRequest
	data := ctx.Body()

	if err := json.Unmarshal(data, &createTag); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(createTag); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createTag", createTag)
	return ctx.Next()
}

func CampanhaValidationMiddleware(ctx *fiber.Ctx) error {
	zap.L().Info("Starting campanha validation")

	var createCampanha dtos.CreateCampanhaRequest
	data := ctx.Body()

	if err := json.Unmarshal(data, &createCampanha); err != nil {
		zap.L().Error("Error when unmarshalling data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid field type",
		})
	}

	if err := Validate.Struct(createCampanha); err != nil {
		var jsonValidationError validator.ValidationErrors
		if errors.As(err, &jsonValidationError) {
			errorsCauses := []exceptions.Causes{}
			for _, e := range jsonValidationError {
				cause := exceptions.Causes{
					FieldMessage: e.Translate(transl),
					Field:        e.Field(),
				}
				errorsCauses = append(errorsCauses, cause)
			}
			zap.L().Error("Error validating fields", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"request invalid": exceptions.NewBadRequestValidationError("Some fields are invalid", errorsCauses),
			})
		}

		zap.L().Info("Error converting fields", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error trying to convert fields",
		})
	}

	ctx.Locals("createCampanha", createCampanha)
	return ctx.Next()
}
