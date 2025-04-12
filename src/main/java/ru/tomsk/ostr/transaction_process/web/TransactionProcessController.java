package ru.tomsk.ostr.transaction_process.web;

import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import ru.tomsk.ostr.transaction_process.dto.PaymentDto;
import ru.tomsk.ostr.transaction_process.service.TransactionProcessService;

import java.util.List;

@RestController
@RequiredArgsConstructor
public class TransactionProcessController {
    private final TransactionProcessService service;

    @GetMapping("/api/transactions")
    public List<PaymentDto> getLast(@RequestParam(name = "count") final Integer count) {
        return service.getLast(count);
    }
}
