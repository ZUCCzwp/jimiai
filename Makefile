## 构建app	
.PHONY: buildApp
buildApp:
	go build -o aiApp ./cmd/app/main.go

## 运行app
.PHONY: runApp
runApp:
	./aiApp

## 构建并运行app	
.PHONY: buildAndRunApp
buildAndRunApp:
	go build -o aiApp ./cmd/app/main.go && ./aiApp

## 构建并运行admin	
.PHONY: buildAndRunAdmin
buildAndRunAdmin:
	go build -o aiAdmin ./cmd/admin/main.go && ./aiAdmin
