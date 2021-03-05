run2:
	@export  BIND_ADDRESS=:8080; \
	export   APP_SESSION_NAME="WAUTH-ALTOFS"; \
	export   APP_SESSION_KEY="wauth-altofs-key"; \
	export   AUTH_GOOGLE_CLIENT_ID="572098240064-j1fprk3744u8g6vhrejaclicpm595irm.apps.googleusercontent.com"; \
	export   AUTH_GOOGLE_CLIENT_SECRET="zZ1jrVLW3A7-chpRnfoQvaOS"; \
	export   AUTH_GOOGLE_REDIRECT_URL="http://localhost:8080/auth/google/callback"; \
	go run main.go

run:
	go run main.go

run3:
	@export WAUTH_LOG_LEVEL="debug"; \
	 export WAUTH_SFTP_AUTH_CHECK_IP=0; \
	 export WAUTH_BIND_ADDRESS=":8080"; \
	 export WAUTH_APP_URL="http://localhost:8080"; \
	 export WAUTH_AUTH_GOOGLE_CLIENT_ID="922970115581-jt213kv85j0bab41ut9si1u1d3saq2f2.apps.googleusercontent.com"; \
	 export WAUTH_AUTH_GOOGLE_CLIENT_SECRET="2qrgJ9L4uKB7Mt6g3ACDl3qb"; \
     export WAUTH_AUTH_GOOGLE_REDIRECT_URL="http://localhost:8080/auth/google/callback"; \
	 export WAUTH_SFTP_SERVER="localhost:2322"; \
     export WAUTH_SFTP_SERVER_URL="http://127.0.0.1:7070"; \
	 export WAUTH_SFTP_SERVER_URL_USER="api"; \
	 export WAUTH_SFTP_SERVER_URL_PASS="Api@p!"; \
	go run main.go
