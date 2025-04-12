package ru.tomsk.ostr.transaction_process.dto;

import lombok.Builder;
import lombok.extern.jackson.Jacksonized;
import ru.tomsk.ostr.transaction_process.domain.PaymentEntity;

@Builder
@Jacksonized
public record PaymentDto(String id, Long paymentDate, String sender, String receiver, Float amount) {
    public static PaymentDto of(final PaymentEntity entity) {
        return PaymentDto.builder()
                .id(entity.getId())
                .paymentDate(entity.getPaymentDate().getTime())
                .sender(entity.getSender())
                .receiver(entity.getReceiver())
                .amount(entity.getAmount())
                .build();
    }
}
