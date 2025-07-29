CREATE TABLE fornecedores (
    id INT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    telefone VARCHAR(20),
    email VARCHAR(150),
    endereco VARCHAR(255),
    cidade VARCHAR(100),
    estado VARCHAR(2),
    cep VARCHAR(10),
    data_cadastro DATE,
    status VARCHAR(20)
);

CREATE TABLE produtos (
    id INT PRIMARY KEY,               -- ID do produto (chave primária)
    codigo_barra VARCHAR(20) NOT NULL, -- Código de barras do produto
    nome_produto VARCHAR(255) NOT NULL, -- Nome do produto
    sku VARCHAR(50) NOT NULL,          -- SKU do produto
    categoria VARCHAR(100) NOT NULL,   -- Categoria do produto
    destinado_para VARCHAR(100) NOT NULL, -- Destinado para (ex: Cães, Gatos, etc.)
    variacao VARCHAR(100),            -- Variação do produto (ex: Raças Grandes, Filhotes, etc.)
    marca VARCHAR(100),               -- Marca do produto
    descricao TEXT,                   -- Descrição do produto
    status VARCHAR(20) NOT NULL,      -- Status do produto (ativo, inativo, etc.)
    preco_venda DECIMAL(10, 2) NOT NULL -- Preço de venda do produto
);

CREATE TABLE pedido_compra (
    id_pedido INT PRIMARY KEY,                -- ID do pedido de compra
    id_fornecedor INT,                        -- ID do fornecedor (referência à tabela fornecedores)
    data_pedido DATE NOT NULL,                -- Data do pedido
    data_entrega DATE NOT NULL,               -- Data de entrega
    valor_frete DECIMAL(10, 2),               -- Valor do frete
    custo_pedido DECIMAL(10, 2),              -- Custo do pedido
    valor_total DECIMAL(10, 2),               -- Valor total (custo pedido + frete)
    descricao_pedido TEXT,                    -- Descrição do pedido
    status VARCHAR(20),                       -- Status do pedido (entregue, em transporte, etc.)
    FOREIGN KEY (id_fornecedor) REFERENCES fornecedores(id) -- Relacionamento com a tabela fornecedores
);

CREATE TABLE item_pedido (
    id_item INT PRIMARY KEY,                 -- ID do item do pedido
    id_pedido INT,                           -- ID do pedido (referência à tabela `pedido_compra`)
    id_produto INT,                          -- ID do produto (referência à tabela `produtos`)
    quantidade INT NOT NULL,                 -- Quantidade do produto no pedido
    preco_unitario DECIMAL(10, 2) NOT NULL,  -- Preço unitário do produto
    subtotal DECIMAL(10, 2) NOT NULL,        -- Subtotal do item (quantidade * preco_unitario)
    FOREIGN KEY (id_pedido) REFERENCES pedido_compra(id_pedido),  -- Relacionamento com a tabela `pedido_compra`
    FOREIGN KEY (id_produto) REFERENCES produtos(id)               -- Relacionamento com a tabela `produtos`
);

CREATE TABLE produto_precos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,         -- ID do registro
    id_produto INT,                            -- ID do produto (referência à tabela produtos)
    preco_venda DECIMAL(10, 2) NOT NULL,       -- Preço de venda do produto
    cmv DECIMAL(10, 2) NOT NULL,               -- Custo das mercadorias vendidas (CMV)
    margem DECIMAL(5, 2) NOT NULL,             -- Margem de lucro calculada (preço_venda - cmv) / preco_venda
    data_registro DATE NOT NULL,               -- Data do registro da alteração de preço ou CMV
    FOREIGN KEY (id_produto) REFERENCES produtos(id)  -- Relacionamento com a tabela produtos
);
