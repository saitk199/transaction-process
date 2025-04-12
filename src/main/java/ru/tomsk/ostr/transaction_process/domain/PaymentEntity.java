package ru.tomsk.ostr.transaction_process.domain;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.*;

import java.sql.Timestamp;

@Setter
@Getter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity(name = "Payment")
public class PaymentEntity {

    @Id
    private String id;

    private Timestamp paymentDate;

    private String sender;

    private String receiver;

    private Float amount;

}
