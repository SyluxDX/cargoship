FTPUSER=${FTPUSER:-ftpuser}
IP=${IP:-0.0.0.0}

# run server
docker run -d -v $(pwd)/ftpdata:/home/vsftpd/$FTPUSER \
-p 20:20 -p 21:21 -p 21100-21110:21100-21110 \
-e FTP_USER=$FTPUSER -e FTP_PASS=passwd \
-e PASV_ADDRESS=$IP -e PASV_MIN_PORT=21100 -e PASV_MAX_PORT=21110 \
--name vsftpd --rm fauria/vsftpd
