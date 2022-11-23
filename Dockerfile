FROM debian


ADD bin/ /bin
COPY configs/app.env /app.env

EXPOSE 8080

ENTRYPOINT ["/bin/control-agent"]
