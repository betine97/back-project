/*

INSERT INTO produtos (codigo_barra, nome_produto, sku, categoria, destinado_para, variacao, marca, descricao, status, preco_venda)
VALUES
('7891234560001', 'Ração 15 kg Premium Cães Adultos', 'RAC-PREM-CAES-ADU', 'Alimentação', 'Cães', 'Raças Grandes', 'Golden', 'Ração completa para cães adultos de grande porte.', 'ativo', 189.90),
('7891234560002', 'Ração 10 kg Premium Gatos Filhotes', 'RAC-PREM-GAT-FIL', 'Alimentação', 'Gatos', 'Filhotes', 'Whiskas', 'Ração seca especialmente formulada para gatos filhotes.', 'ativo', 99.90);


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
*/