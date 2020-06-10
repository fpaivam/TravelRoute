# TravelRoute

Um turista deseja viajar pelo mundo pagando o menor preço possível, independentemente do número de conexões necessárias.
Vamos construir um programa que facilite ao nosso turista escolher a melhor rota para sua viagem.

## Exemplo de entrada
```csv
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
```

## Compilar

Para compilar este programa basta:
```bash
go build
```

## Testes

Para rodar todos os testes basta:

```bash
go test -p 1 ./...
```

É necessário rodar os testes sem paralelização, pois os testes do webserver quando rodam simultaneamente tentam abrir uma mesma porta causando falsos negativos

## Rodando o programa

Para rodar o programa basta:

```bash
./TravelRoute providedInput.csv
```

## Estrutura dos pacotes

Este programa contém 5 pacotes:
- main
- algorithm
- controller
- dal
- domain

_main_ é o pacote que gera o executável (onde se encontra a função main). Este pacote é responsável por decodificar os argumentos da linha de comando e inicializar algumas estruturas

_algorithm_ contem o código responsável por gerenciar o webserver HTTP e suas rotas

_controller_ contem o código responsavel por genrenciar o webserver HTTP e suas rotas

_dal_ contém toda a lógica de acesso aos dados, neste caso o arquivo CSV

_domain_ contém toda a lógica de negócio do programa. Responsável por encontrar a rota mais barata.

## API REST

Este programa contém 2 endpoints:
- _/route_
- _/route/best_

### /route

É responsável por gerir os dados das rotas. Aceita GET, POST e PUT

#### GET /route

Lista as rotas cadastradas. Exemplo de retorno:

```json
[
    {
        "Origin": "GRU",
        "Destination": "BRC",
        "Cost": 10
    },
    {
        "Origin": "BRC",
        "Destination": "SCL",
        "Cost": 5
    },
    {
        "Origin": "GRU",
        "Destination": "CDG",
        "Cost": 75
    },
    {
        "Origin": "GRU",
        "Destination": "SCL",
        "Cost": 20
    },
    {
        "Origin": "GRU",
        "Destination": "ORL",
        "Cost": 56
    },
    {
        "Origin": "ORL",
        "Destination": "CDG",
        "Cost": 5
    },
    {
        "Origin": "SCL",
        "Destination": "ORL",
        "Cost": 20
    }
]
```

#### POST /route

Insere uma nova rota. Exemplo de Envio:
```json
{
    "Origin": "SCL",
    "Destination": "GRU",
    "Cost": 2
}
```
Exemplo de retorno:
```
OK
```

### /route/best

É responsável por encontrar a rota mais barata entre _Origin_ e _Destination_. Aceita somente GET.

#### GET /route/best

Procura a rota mais barata entre _Origin_ e _Destination_. Estes parâmetros devem ser passados na query string. Exemplo:

Get /route/best?Origin=GRU&Destination=CDG
```json
{
    "Route": [
        "GRU",
        "BRC",
        "SCL",
        "ORL",
        "CDG"
    ],
    "Cost": 40
}
```