buildBackend:
	docker compose up --build

startBackend:
	docker compose up

startNoLogsBackend:
	docker compose up -d


removeContainer:
	docker compose down


removeContainerData:
	docker compose down --volumes


startApp:
	docker compose up -d
	cd frontend && npm run dev





.PHONY: build start startNoLogs removeContainer removeContainerData startApp
