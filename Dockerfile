FROM debian:9-slim
WORKDIR /var/project

ARG uid=1000
ARG user=pgmycli

RUN apt-get update && apt-get install -y wget gnupg2 lsb-release \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list \
    && echo "deb http://repo.mysql.com/apt/debian $(lsb_release -sc) mysql-8.0" | tee /etc/apt/sources.list.d/mysql80.list \
    && apt-get update && apt-get install -y --allow-unauthenticated \
    postgresql-client-13 \
    mysql-client \
    && apt-get clean \
    && apt-get autoclean\
    && rm -rf /var/lib/apt/lists/* \
    && useradd -u $uid -d /home/$user $user \
    && mkdir /home/$user \
    && chown $user:$user /home/$user

USER $user
