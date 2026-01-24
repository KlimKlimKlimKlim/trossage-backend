#!/bin/zsh

export PATH="$HOME/.cargo/bin:$PATH"

PROJECT_DIR="$1"
FILE="$2"

if [[ -z "$PROJECT_DIR" || -z "$FILE" ]]; then
    echo "❌ Error: Project directory or file not provided"
    echo "Usage: $0 /path/to/project /path/to/file"
    exit 1
fi

cd "$PROJECT_DIR" || exit 1

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🔍 Linting: $(basename "$FILE")"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [[ "$FILE" == *.go ]]; then
    echo "🐹 golangci-lint"
    FILE_DIR=$(dirname "$FILE")

    golangci-lint run --config=.golangci.yml --fix "$FILE_DIR"
    EXIT=$?

elif [[ "$FILE" == *.sql ]]; then
    echo "🔒 squawk"
    squawk --config=.squawk.toml "$FILE"
    SQUAWK_EXIT=$?

    echo ""
    echo "──────────────────────────────────────────"
    echo ""

    echo "🐘 sqlfluff"
    sqlfluff lint --config=.sqlfluff "$FILE"
    SQLFLUFF_EXIT=$?

    EXIT=$((SQUAWK_EXIT + SQLFLUFF_EXIT))

else
    echo "⚠️  Unsupported file type: $(basename "$FILE")"
    echo "Supported: .go, .sql"
    echo ""
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
[ $EXIT -eq 0 ] && echo "  ✅ Pass" || echo "  ❌ Fail"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

exit $EXIT
