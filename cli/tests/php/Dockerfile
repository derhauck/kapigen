FROM php:8.2
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer
RUN apt update && apt install -y git && mkdir /.composer && chown 1000:1000 /.composer
RUN pecl install xdebug && docker-php-ext-enable xdebug