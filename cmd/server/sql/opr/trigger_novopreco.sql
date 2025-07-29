-- Trigger para SQLite - quando preço de venda muda
CREATE TRIGGER IF NOT EXISTS after_preco_venda_update
AFTER UPDATE ON produtos
FOR EACH ROW
WHEN NEW.preco_venda != OLD.preco_venda
BEGIN
    INSERT INTO produto_precos (id_produto, preco_venda, cmv, margem, data_registro)
    SELECT 
        NEW.id,
        NEW.preco_venda,
        -- Assumindo CMV como 60% do preço de venda (ajuste conforme necessário)
        NEW.preco_venda * 0.6,
        ROUND(((NEW.preco_venda - (NEW.preco_venda * 0.6)) / NEW.preco_venda) * 100, 2),
        DATE('now')
    WHERE NOT EXISTS (
        SELECT 1 FROM produto_precos 
        WHERE id_produto = NEW.id 
        AND preco_venda = NEW.preco_venda 
        AND DATE(data_registro) = DATE('now')
    );
END;