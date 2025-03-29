-- Заполняем таблицу vallet
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr1', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr2', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr3', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr4', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr5', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr6', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr7', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr8', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr9', 100.00);
INSERT OR IGNORE INTO vallet (adress, balance) VALUES ('addr10', 100.00);

-- Заполняем таблицу payment
INSERT OR IGNORE INTO payment (id, "sender", "recipient", payment_date, amount)
    VALUES ('tx1', 'addr1', 'addr2', 1712000000, 10.0);
INSERT OR IGNORE INTO payment (id, "sender", "recipient", payment_date, amount)
    VALUES ('tx2', 'addr2', 'addr1', 1712100000, 15.0);
