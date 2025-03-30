### Requirements

| package        | version   | web-page                                                          |
|----------------|-----------|:-----------------------------------------------------------------:|
| Ubuntu         | 24.04     | [🌐](https://ubuntu.com/download)                                 |
| Docker         | 27.5.1-rd | [🌐](https://docs.docker.com/desktop/setup/install/linux/ubuntu/) |
| Docker-Compose | 2.33.0    | [🌐](https://docs.docker.com/desktop/setup/install/linux/ubuntu/) |
| golang         | 1.24.1    | [🌐](https://go.dev/)                                             |
| air-verse/air  | 1.67.7    | [🌐](https://github.com/air-verse/air)                            |
| direnv         | 2.35.0    | [🌐](https://github.com/direnv/direnv?tab=readme-ov-file)         |
| golang-migrate | 4.18.2    | [🌐](https://github.com/golang-migrate/migrate)                   |

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