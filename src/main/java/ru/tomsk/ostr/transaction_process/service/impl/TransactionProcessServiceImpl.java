package ru.tomsk.ostr.transaction_process.service.impl;

import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import ru.tomsk.ostr.transaction_process.domain.PaymentEntity;
import ru.tomsk.ostr.transaction_process.dto.PaymentDto;
import ru.tomsk.ostr.transaction_process.dto.WalletDto;
import ru.tomsk.ostr.transaction_process.repository.PaymentRepository;
import ru.tomsk.ostr.transaction_process.repository.WalletRepository;
import ru.tomsk.ostr.transaction_process.service.TransactionProcessService;

import java.sql.Timestamp;
import java.time.Instant;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class TransactionProcessServiceImpl implements TransactionProcessService {
    private final PaymentRepository paymentRepository;
    private final WalletRepository walletRepository;

    @Override
    public PaymentDto send(final PaymentDto payment) {
        var sender = getBalance(payment.sender());
        if (sender.balance() < payment.amount()) {
            throw new RuntimeException("Не достаточно средств");
        }
        var receiver = getBalance(payment.receiver());

        saveBalance(sender.address(), sender.balance() - payment.amount());

        saveBalance(receiver.address(), receiver.balance() + payment.amount());

        var saved = PaymentEntity.builder()
                .id(UUID.randomUUID().toString())
                .paymentDate(Timestamp.from(Instant.now()))
                .sender(payment.sender())
                .receiver(payment.receiver())
                .amount(payment.amount())
                .build();

        paymentRepository.save(saved);

        return PaymentDto.of(saved);
    }

    private WalletDto saveBalance(String address, Float newBalance) {
        var wallet = walletRepository.findById(address)
                .orElseThrow(() -> new RuntimeException("Ошибка при получении кошелька по адресу: " + address));
        wallet.setBalance(newBalance);
        var saved = walletRepository.save(wallet);
        return WalletDto.of(saved);
    }

    @Override
    public WalletDto getBalance(final String address) {
        return WalletDto.of(walletRepository.findById(address)
                .orElseThrow(() -> new RuntimeException("Ошибка при получении кошелька по адресу: " + address)));
    }

    @Override
    public List<PaymentDto> getLast(final Integer count) {
        Sort sort = Sort.by(Sort.Order.by("paymentDate"));
        Pageable limit = PageRequest.of(0, count, sort);
        return paymentRepository.findAll(limit).stream()
                .map(PaymentDto::of)
                .toList();
    }
}
