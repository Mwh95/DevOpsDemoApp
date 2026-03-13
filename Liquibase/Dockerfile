FROM liquibase/liquibase:4.20.0

COPY ./modules/database/bootstrap/ /liquibase/modules/database/bootstrap/
COPY ./modules/mapservice/changelog/ /liquibase/modules/mapservice/changelog/
COPY ./run-migrations.sh /liquibase/run-migrations.sh

ENTRYPOINT ["sh", "/liquibase/run-migrations.sh"]
