package ru.tomsk.ostr.transaction_process.service;

import org.springframework.lang.NonNull;
import org.springframework.lang.Nullable;
import ru.tomsk.ostr.transaction_process.dto.PaymentDto;
import ru.tomsk.ostr.transaction_process.dto.WalletDto;

import java.util.List;

public interface TransactionProcessService {
    PaymentDto send(@NonNull PaymentDto payment);
    WalletDto getBalance(@NonNull String address);
    List<PaymentDto> getLast(@Nullable Integer count);
}
