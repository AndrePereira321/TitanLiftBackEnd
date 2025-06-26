@echo off
echo 🧪 Running tests before push...

REM Run all Go tests

go test ./...

IF %ERRORLEVEL% EQU 0 (
    echo ✅ Tests passed! Proceeding with push.
    EXIT /B 0
) ELSE (
    echo ❌ Tests failed! Push cancelled.
    EXIT /B 1
)