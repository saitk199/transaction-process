package ru.tomsk.ostr.transaction_process.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import ru.tomsk.ostr.transaction_process.domain.WalletEntity;

public interface WalletRepository extends JpaRepository<WalletEntity, String> {
}
