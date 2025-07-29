-- Trigger para SQLite - quando item do pedido é inserido
CREATE TRIGGER IF NOT EXISTS after_item_pedido_insert
AFTER INSERT ON item_pedido
FOR EACH ROW
BEGIN
    INSERT INTO produto_precos (id_produto, preco_venda, cmv, margem, data_registro)
    SELECT 
        NEW.id_produto,
        p.preco_venda,
        NEW.preco_unitario, -- Usando preço unitário como CMV
        ROUND(((p.preco_venda - NEW.preco_unitario) / p.preco_venda) * 100, 2),
        DATE('now')
    FROM produtos p
    WHERE p.id = NEW.id_produto
    AND NOT EXISTS (
        SELECT 1 FROM produto_precos pp
        WHERE pp.id_produto = NEW.id_produto 
        AND pp.cmv = NEW.preco_unitario
        AND DATE(pp.data_registro) = DATE('now')
    );
END;