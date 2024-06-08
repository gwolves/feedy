env "local" {
  src = "file://database/schema.sql"
  url = "postgres://postgres:pass@localhost:5432/feedy?sslmode=disable"
  dev = "docker://postgres/16/dev?search_path=public"

  migration {
      dir = "file://database/migrations"
  }
}
