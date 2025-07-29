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



INSERT INTO pedidos (
    id_fornecedor, data_pedido, data_entrega, valor_frete, custo_pedido, valor_total, descricao_pedido, status
) VALUES
(1, '2024-08-01', '2024-08-05', 100.00, 2000.00, 2100.00, 'Compra de ração e acessórios para cães.', 'entregue'),
(2, '2024-08-02', '2024-08-06', 150.00, 3000.00, 3150.00, 'Compra de produtos para gatos e aves.', 'em transporte');

INSERT INTO item_pedido (
    id_item, id_pedido, id_produto, quantidade, preco_unitario, subtotal
) VALUES
(1, 1, 1, 30, 100.00, 3000.00),
(2, 1, 4, 40, 10.00, 400.00);


INSERT INTO produto_precos (id_produto, preco_venda, cmv, margem, data_registro)
SELECT 
    ip.id_produto,                                 -- Captura o id_produto da tabela item_pedido
    p.preco_venda,                                 -- Captura o preco_venda da tabela produtos
    SUM(ip.preco_unitario * ip.quantidade) AS cmv,  -- Calcula o cmv (custo das mercadorias vendidas)
    (p.preco_venda - SUM(ip.preco_unitario * ip.quantidade)) / p.preco_venda AS margem,  -- Calcula a margem
    ped.data_pedido                                -- Data de registro capturada da tabela pedidos
FROM item_pedido ip
JOIN produtos p ON ip.id_produto = p.id          -- Relacionamento com a tabela produtos
JOIN pedidos ped ON ip.id_pedido = ped.id_pedido -- Relacionamento com a tabela pedidos
GROUP BY ip.id_produto, p.preco_venda, ped.data_pedido;
