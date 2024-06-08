default:
    @just --list

build:
    @go build -o build/feedy

run *args:
    @go run . {{args}}

up:
    @docker-compose up -d db
    @go run . runserver

makemigration *args:
  atlas migrate diff --env local {{args}}

migrate:
  atlas migrate apply --env local
