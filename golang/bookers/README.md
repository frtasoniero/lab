### Requirements

| package        | version   | web-page                                                     |
|----------------|-----------|:------------------------------------------------------------:|
| golang         | v1.24.1   | [🌐](https://go.dev/)                                        |
| chi            | v5.2.1    | [🌐](https://github.com/go-chi/chi)                          |
| air-verse/air  | v1.67.7   | [🌐](https://github.com/air-verse/air)                       |
| direnv         | v.2.35    | [🌐](https://github.com/direnv/direnv?tab=readme-ov-file)    |
| lib/pg         | v1.10.9   | [🌐](https://github.com/lib/pq)                              |
| golang-migrate | v4.18.2   | [🌐](https://github.com/golang-migrate/migrate)              |

### Migrations

- CREATE USER
> ```
> make migration create_user
> ```

- MIGRATE UP
> ```
> make migrate-up
> ```

- MIGRATE DOWN
> ```
> make migrate-down
> ```