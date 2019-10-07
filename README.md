# e-Commerce

API para el consumo de plataforma e-Commerce con pago mediante criptomonedas

## Requisitos
- Go v1.12 como mínimo
- Postman

## ¿Cómo correr?

Go get del repositorio

```
go get github.com/jcamilom/ecommerce
```

Navegar a la carpeta dentro del GOPATH, normalmente el directorio HOME

```
cd ~/go/src/github.com/jcamilom/ecommerce
```

Setear los valores `AWS_ACCESS_KEY_ID` y `AWS_SECRET_ACCESS_KEY` en el archivo `.env`

Correr el programa

```
go run main.go
```

Utilizar colección de postman para testear las diferentes funcionalidades.

## Arquitectura

- Pago criptomendas: [Stellar Network](https://www.stellar.org/)
- Persistencia de datos: Amazon DynamoDB (NoSQL)
- Servicio monolith style en Golang con /gorilla/mux

### Stellar Network

Los pagos en la plataforma se hacen por medio del testnet de la red Stellar. Al momento de creación de cuentas en la plataforma e-Commerce se registra un keypair en el testnet de Stellar y se carga la dirección con 10.000 lumens listos para ejecutar transacciones en la red.

La plataforma e-Commerce cuenta con una dirección en el testnet registrada donde se hacen las transferencias de las compras de los usuarios.

### Persistencia de datos

El API está respaldado por tres bases de datos en DynamoDB: `Users`, `Purchases` y `Products`.

---
#### User model (Table)

Representa un usuario registrado en la plataforma

| Field         | Type          |
| ------------- |:-------------:|
| ID      | string |
| Name      | string    |
| Email | string      |
| PasswordHash | string      |
| AccessToken | string      |
| Favorites | []Favorite     |
| Wallet | Wallet     |

#### Favorite model
| Field         | Type          |
| ------------- |:-------------:|
| ID      | string |
| Name      | string    |
| Price | number      |

#### Wallet model
| Field         | Type          |
| ------------- |:-------------:|
| Seed      | string |
| Address      | string    |

---
#### Purchases model (Table)

Representa un registro de item comprado en la plataforma por algún usuario

| Field         | Type          |
| ------------- |:-------------:|
| ID      | string |
| Name      | string    |
| Date | string      |
| Item | PurchaseItem      |

#### PurchaseItem model
| Field         | Type          |
| ------------- |:-------------:|
| ID      | string |
| Name      | string    |
| Price | number      |

---

Representa un producto del catálogo de la plataforma

Existen 5 productos en el catálogo catálogo de la tienda, con id del 1-5. El id del item es necesario en el body de la petición de compra de ítems.

#### Products model (Table)
| Field         | Type          |
| ------------- |:-------------:|
| ID      | string |
| Name      | string    |
| Price | number      |
| Quantity | number      |

---

## Faltantes del entregable

- Pruebas unitarias
- Funcionalidad carrito de compras
