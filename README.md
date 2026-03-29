<br />
<div align="center">
  <a href="https://github.com/github_username/repo_name">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Fullstack Hotell App</h3>

## About The Project

Репозиторий с готовым решением для ведения собсвенного отеля. Проект содержит главную страницу, внутренний бэкенд бронирования, а так же админ панель для управления бронированиями.

<!-- GETTING STARTED -->
## Getting Started

На данном этапе возмежен локальный запуск проекта. (docker image in progress)
Скачайте проект используя команду `git clone https://github.com/emptywe/local_hotel`
Для начала работы потребуется устанвка go не ниже версии 1.20, Postgresql 12 и выше. Установить их можно пользуясь инструкцией с официальных сайтов проектов.
После устанвки go в репозитроии потребуется выполнить команду go mod tidy для устанволения всех зависимостей. 


```sh
go mod tidy
```

После установки postgresql требуется создать соответствующую проекту БД - bookings. после запуска проекта все нужные миграции автоматически встанут в базу.

```sql
CREATE DATABASE bookings;
```

Для запуска проекта воспользуйтсь командой go run `go run cmd/web/main.go cmd/web/send_mail.go`.
На localhost:8080 вы сможете увдиеть главную страницу отеля.