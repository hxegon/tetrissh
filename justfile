_default:
  @just --list

run:
  go run ./cmd/tetrissh/main.go

conn:
  ssh ssh://localhost:42069

watch:
  fd -e go | entr -cr go run ./cmd/tetrissh/main.go

watch-conn:
  while true; do just conn; sleep 2; done
