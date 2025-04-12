package ru.tomsk.ostr.transaction_process.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.*;

@Setter
@Getter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity(name = "WALLET")
public class WalletEntity {
    @Id
    @Column(name = "ADDRESS")
    String address;

    @Column(name = "BALANCE")
    Float balance;
}
