# go-url-short

## TODO
- [ ] Deploy to AWS Lambda 
- [ ] Add tests
- [ ] Add CI/CD
- [x] 표준 프로젝트 구조로 수정
- [x] postgresql 연결 및 store 구성
- [ ] 62 진법 적용 
- [ ] 운영 배포 (도메인, 서버, DB)

```shell
docker run --rm \
    --name postgresql \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -e POSTGRES_USER="username" \
    -e POSTGRES_DB=shorturl \
     -v ./db/db.sql:/docker-entrypoint-initdb.d \
    postgres
```
