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