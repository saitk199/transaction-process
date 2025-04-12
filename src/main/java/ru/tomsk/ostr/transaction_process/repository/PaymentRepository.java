package ru.tomsk.ostr.transaction_process.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import ru.tomsk.ostr.transaction_process.domain.PaymentEntity;

public interface PaymentRepository extends JpaRepository<PaymentEntity, String> {
}
