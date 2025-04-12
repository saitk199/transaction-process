package ru.tomsk.ostr.transaction_process.dto;

import lombok.Builder;
import lombok.extern.jackson.Jacksonized;
import ru.tomsk.ostr.transaction_process.domain.WalletEntity;

@Builder
@Jacksonized
public record WalletDto(String address, Float balance) {
    public static WalletDto of(WalletEntity entity) {
        return WalletDto.builder()
                .address(entity.getAddress())
                .balance(entity.getBalance())
                .build();
    }
}
