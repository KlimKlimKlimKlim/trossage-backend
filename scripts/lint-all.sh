#!/bin/zsh

export PATH="$HOME/.cargo/bin:$PATH"

PROJECT_DIR="$1"

if [[ -z "$PROJECT_DIR" ]]; then
    echo "❌ Error: Project directory not provided"
    echo "Usage: $0 /path/to/project"
    exit 1
fi

cd "$PROJECT_DIR" || exit 1

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🔍 Running linters"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Go
echo "🐹 golangci-lint (Go)"
golangci-lint run --config=.golangci.yml
GO_EXIT=$?
[ $GO_EXIT -eq 0 ] && echo "✅ Pass" || echo "❌ Fail"

echo ""
echo "──────────────────────────────────────────"
echo ""

# SQL Safety
echo "🔒 squawk (SQL safety)"
squawk --config=.squawk.toml migrations/*.sql internal/repository/postgres/query/*/*.sql
SQUAWK_EXIT=$?
[ $SQUAWK_EXIT -eq 0 ] && echo "✅ Pass" || echo "❌ Fail"

echo ""
echo "──────────────────────────────────────────"
echo ""

# SQL Style
echo "🐘 sqlfluff (PostgreSQL style)"
sqlfluff lint --config=.sqlfluff migrations/*.sql internal/repository/postgres/query/*/*.sql
SQLFLUFF_EXIT=$?
[ $SQLFLUFF_EXIT -eq 0 ] && echo "✅ Pass" || echo "❌ Fail"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Exit
TOTAL=$((GO_EXIT + SQUAWK_EXIT + SQLFLUFF_EXIT))
if [ $TOTAL -eq 0 ]; then
    echo "  ✅ All checks passed"
else
    echo "  ❌ Some checks failed"
fi
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

exit $TOTAL
