# wildcard-processor
Wildcard-processor adds wildcard support for recipient hosts validation for 
[go-guerilla](https://github.com/flashmob/go-guerrilla) package

## About

This package is a _Processor_ for the [Go-Guerrilla SMTP server](https://github.com/flashmob/go-guerrilla). By default it is possible to
match hosts using exact names or "." character. Wildcard-processor adds another configuration option,
where it is possible to define recipient hosts using wildcard, e.g. "*.com", so it offers much greater
flexibility.

## Configuration

Set `wildcard` as _validate_process_ in your backend config file and define hosts with 
wildcards under _wildcard_hosts_ configuration field. Use commas for multiple values.

```json
"backend_config":
{
    "validate_process": "wildcard",
    "wildcard_hosts": "*.com,*.org",
    "log_received_mails": false,
},
```

Then import `github.com/DevelHell/wildcard-processor` and add wildcard as a processor

```go
app.AddProcessor("wildcard", wildcard_processor.WildcardProcessor)
```

And you're ready to go. For more information see [go-guerilla documentation](https://github.com/flashmob/go-guerrilla/wiki/Using-as-a-package#registering-a-processor)
