
# impersonate a service account using golang

- We are authenticated using `impersonating-service-account` 
- `impersonating-service-account` is granted role `Service Account Token Creator` permissions on a service account `target-service-account`
- `target-service-account` has permission to read from a bucket 
- `impersonating-service-account` has no  permission to read from this bucket
- `impersonating-service-account` impersonates the `target-service-account` to enumerate the contents of a bucket.

To test locally :

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/impersonating-service-account-key.json 
export TARGET_SERVICE_ACCOUNT=target-service-account@yourproject.iam.gserviceaccount.com 
go run ./main.go
```