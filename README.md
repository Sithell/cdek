# cdek
Go client library for CDEK REST API

## Usage

```go
package main

import (
	"fmt"
	"github.com/Sithell/cdek"
)

func main() {
	client, err := cdek.NewClientWithBaseUrl(
		"EMscd6r9JnFiQ3bLoyjJY6eM78JrJceI",
		"PjLZkKBHEiLK3YsjtNrt3TGNG0ahs3kG",
		cdek.BaseUrlV2Test,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	tariffCosts, err := client.GetShippingCost(
		"Россия, г. Москва, Cлавянский бульвар д.1",
		"Россия, Воронежская обл., г. Воронеж, ул. Ленина д.43",
		[]cdek.Package{{10, 10, 40, 3000}},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tariffCosts)
}
```
