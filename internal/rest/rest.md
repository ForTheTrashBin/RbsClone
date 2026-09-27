1. Normalfall (Erfolgreiche Antwort)

1.1 Der Huma-Standart
1.1.1 Bei "return &ResponseStruct(), nil" --> Status: 200
1.1.2 Bei "return nil, nil" --> Status: 204
Der Erfolgscode kann einheitlich für beide Antwortarten gemeinsam (1.1.1 und 1.1.2) mit 'Defaultstatus' auf einen anderen gewünschten Wert gesetzt werden.

```go
huma.Register(api, huma.Operation{
    Method:        http.MethodPost,
    Path:          "/users",
    DefaultStatus: 201, // <-- Hier wird der Erfolgscode vordefiniert
}, func(ctx context.Context, input *MyInput) (*MyOutput, error) {
    return &MyOutput{...}, nil
})
```

2. Fehlerfall ("return nil, err")

2.1 Vorgefertigte Huma-Fehler
    "return nil, huma.Error404NotFound" oder "return nil, huma.Error400BadRequest"

2.2 Standard-Go-Fehler (z.B.: fmt.Errorf("..."))
    Diese Fehler werden als "500 Internal Server Error" umkodiert

3. Automatisch durch die vorgeschaltete Huma-Validierung

3.1 "400 Bad Request", wenn die Request-Daten nicht gelesen werden konnten
3.2 "422 Unprocessable Entity", wenn die inhaltliche Validierung fehlgeschlagen ist



