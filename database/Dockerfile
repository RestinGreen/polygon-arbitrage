FROM postgres:latest



RUN mkdir init
COPY ./pkg/init/ ./init
RUN chmod +x ./init/setup_db.sh

USER postgres 

ENTRYPOINT ["./init/setup_db.sh"]

EXPOSE 5432

CMD ["postgres"]