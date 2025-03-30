### Requirements

>- Golang v1.24.1
>- Chi v5.2.1
>- Air v1.67.7
>- Direnv v.2.35.0
>- Pg v1.10.9
>- Golang-Migrate v4.18.2

### Migrations

CREATE USER
> ```
> $ make migration create_user
> ```

MIGRATE UP
> ```
> $ make migrate-up
> ```

MIGRATE DOWN
> ```
> $ make migrate-down
> ```