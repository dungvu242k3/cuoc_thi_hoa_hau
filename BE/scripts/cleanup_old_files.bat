@echo off
echo === Cleaning up old/duplicate files ===

echo.
echo --- dao/ (6 old _repo files) ---
del "internal\dao\contestant_repo.go"
del "internal\dao\user_repo.go"
del "internal\dao\feedback_repo.go"
del "internal\dao\score_repo.go"
del "internal\dao\schedule_repo.go"
del "internal\dao\local_storage.go"

echo.
echo --- service/ (7 old _svc/_service files) ---
del "internal\service\auth_svc.go"
del "internal\service\auth_service.go"
del "internal\service\contestant_svc.go"
del "internal\service\contestant_service.go"
del "internal\service\feedback_svc.go"
del "internal\service\schedule_svc.go"
del "internal\service\scoring_svc.go"

echo.
echo --- types/ (10 old _port/empty files) ---
del "internal\types\auth_port.go"
del "internal\types\contestant_port.go"
del "internal\types\feedback_port.go"
del "internal\types\score_port.go"
del "internal\types\schedule_port.go"
del "internal\types\security_port.go"
del "internal\types\user_port.go"
del "internal\types\infra.go"
del "internal\types\probe_port.go"
del "internal\types\repository.go"
del "internal\types\service_iface.go"

echo.
echo --- middleware/ (2 old _mdw files) ---
del "graph\middleware\auth_mdw.go"
del "graph\middleware\rate_limit_mdw.go"

echo.
echo --- root cleanup (4 files) ---
del "find_lines.go"
del "build_err.txt"
del "mapper_test_output.txt"
del "docker-compose.yml"

echo.
echo --- configs/ old stubs ---
del "configs\config.go"
del "configs\container.go"

echo.
echo === DONE: 31 files deleted ===
echo Run: go build ./cmd/api/
