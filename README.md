프로젝트 빌드 및 실행 방법
===

Go 언어의 1.11.2번전 환경에서 작성 되었습니다.
Go 패스에 해당 git을 클론시킨 후
cd 
cd $GOPATH/src/wastest/
go build
$GOPATH/wastest
또는
해당 디렉토리에서 go run server.go시행
포트는 기본포트인 8080을 사용합니다.

API가 작동 되고나면  cd $GOPATH/src/wastest/view 디렉토리로 이동후
npm install
npm start 실행 시킨

http://localhost:8081로 접속하면 해당 프로젝트가 작동 되는 것을 볼 수 있습니다.

---

