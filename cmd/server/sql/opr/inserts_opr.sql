INSERT INTO fornecedores (
    id, nome, telefone, email, endereco, cidade, estado, cep, data_cadastro, status
) VALUES
(1, 'Distribuidora Pet Brasil', '(11) 3000-4000', 'contato@petbrasil.com.br', 'Rua dos Animais, 100', 'São Paulo', 'SP', '01000-000', '2024-01-15', 'ativo'),
(2, 'Mega Pet Supply', '(21) 4000-5000', 'vendas@megapetsupply.com.br', 'Av. Pet Lovers, 500', 'Rio de Janeiro', 'RJ', '20000-000', '2024-03-10', 'ativo');


INSERT INTO produtos (id, codigo_barra, nome_produto, sku, categoria, destinado_para, variacao, marca, descricao, status, preco_venda)
VALUES
(1, '7891234560001', 'Ração Premium Cães Adultos', 'RAC-PREM-CAES-ADU', 'Alimentação', 'Cães', 'Raças Grandes', 'Golden', 'Ração completa para cães adultos de grande porte.', 'ativo', 189.90),
(2, '7891234560002', 'Ração Premium Gatos Filhotes', 'RAC-PREM-GAT-FIL', 'Alimentação', 'Gatos', 'Filhotes', 'Whiskas', 'Ração seca especialmente formulada para gatos filhotes.', 'ativo', 99.90),
(3, '7891234560003', 'Suplemento Vitamínico Cães Idosos', 'SUP-VIT-CAES-IDOS', 'Saúde', 'Cães', 'Idosos', 'Organnact', 'Suplemento vitamínico para melhorar a vitalidade de cães idosos.', 'ativo', 54.90),
(4, '7891234560004', 'Brinquedo Bola Mordedor', 'BRI-BOLA-CAES', 'Lazer', 'Cães', 'Raças Pequenas', 'Pet Games', 'Bola mordedora resistente para cães de pequeno porte.', 'ativo', 24.90),
(5, '7891234560005', 'Ração para Pássaros Canário', 'RAC-PASS-CANARIO', 'Alimentação', 'Pássaros', 'Adulto', 'Megazoo', 'Mistura balanceada para canários.', 'ativo', 29.90),
(6, '7891234560006', 'Antipulgas Gatos', 'ANT-PUL-GATOS', 'Saúde', 'Gatos', 'Adulto', 'Bayer', 'Antipulgas em pipeta para gatos adultos.', 'ativo', 89.90),
(7, '7891234560007', 'Arranhador para Gatos', 'ARR-GATOS', 'Lazer', 'Gatos', 'Adulto', 'Chalesco', 'Arranhador compacto para entretenimento e saúde das garras.', 'ativo', 79.90),
(8, '7891234560008', 'Ração Cães Filhotes Raças Pequenas', 'RAC-CAES-FIL-RP', 'Alimentação', 'Cães', 'Filhotes', 'Premier', 'Ração para filhotes de cães de raças pequenas.', 'ativo', 159.90),
(9, '7891234560009', 'Suplemento para Pássaros', 'SUP-PASS-VIT', 'Saúde', 'Pássaros', 'Adulto', 'Avitrin', 'Suplemento vitamínico para todas as espécies de pássaros.', 'ativo', 19.90),
(10, '7891234560010', 'Cordão com Guizo para Gatos', 'COR-GUI-GATOS', 'Lazer', 'Gatos', 'Filhotes', 'Pet Flex', 'Cordão para brincadeira com guizo para gatos filhotes.', 'ativo', 14.90);



INSERT INTO pedido_compra (
    id_pedido, id_fornecedor, data_pedido, data_entrega, valor_frete, custo_pedido, valor_total, descricao_pedido, status
) VALUES
(1, 1, '2024-07-01', '2024-07-05', 150.00, 3850.00, 4000.00, 'Compra mensal de rações, brinquedos e suplementos.', 'entregue'),
(2, 2, '2024-07-01', '2024-07-06', 120.00, 2650.00, 2770.00, 'Compra de produtos para gatos e aves.', 'entregue'),
(3, 1, '2024-07-10', '2024-07-15', 130.00, 1800.00, 1930.00, 'Reposição de brinquedos.', 'em transporte'),
(4, 2, '2024-07-12', '2024-07-17', 140.00, 2100.00, 2240.00, 'Compra emergencial de suplementos.', 'realizado'),
(5, 1, '2024-07-14', '2024-07-20', 160.00, 3200.00, 3360.00, 'Estoque de rações premium.', 'cancelado'),
(6, 2, '2024-07-15', '2024-07-21', 110.00, 2500.00, 2610.00, 'Produtos de higiene animal.', 'recusado');



INSERT INTO item_pedido (
    id_item, id_pedido, id_produto, quantidade, preco_unitario, subtotal
) VALUES
(1, 1, 1, 30, 100.00, 3000.00),
(2, 1, 4, 40, 10.00, 400.00),
(3, 1, 3, 15, 30.00, 450.00),
(4, 2, 2, 25, 55.00, 1375.00),
(5, 2, 5, 20, 15.00, 300.00),
(6, 2, 7, 10, 40.00, 400.00),
(7, 2, 9, 20, 29.00, 580.00),
-- Itens para o pedido 3
(8, 3, 4, 50, 12.00, 600.00),
(9, 3, 7, 30, 40.00, 1200.00),
-- Itens para o pedido 4
(10, 4, 3, 20, 30.00, 600.00),
(11, 4, 6, 25, 50.00, 1250.00),
-- Itens para o pedido 5
(12, 5, 1, 40, 100.00, 4000.00),
(13, 5, 5, 30, 15.00, 450.00),
-- Itens para o pedido 6
(14, 6, 2, 35, 55.00, 1925.00),
(15, 6, 9, 15, 29.00, 435.00);

INSERT INTO produto_precos (id_produto, preco_venda, cmv, margem, data_registro)
VALUES
(1, 150.00, 90.00, 40.00, '2024-01-15'),
(2, 200.00, 120.00, 40.00, '2024-02-10');