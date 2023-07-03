run-scanner:
	@echo "Running scanner..."
	go run cmd/main.go

grab-total:
	@echo "Grabbing total..."
	go run cmd/total.go