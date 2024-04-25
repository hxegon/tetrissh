_default:
  @just

conn:
  ssh ssh://localhost:42069

watch:
  fd -e go | entr -cr go run ./cmd/tetrissh/main.go
