
To test locally :

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/impersonating-service-account-key.json 
export TARGET_SERVICE_ACCOUNT=target-service-account@yourproject.iam.gserviceaccount.com 
go run ./main.go
```