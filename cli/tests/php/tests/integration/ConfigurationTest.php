<?php

namespace Kapigen\Test\Integration;
class ConfigurationTest extends \PHPUnit\Framework\TestCase
{
    function testConfiguration(): void
    {
        $dbname = getenv("DB_NAME");
        $dbhost = getenv("DB_HOST");
        $dbuser = getenv("DB_USER");
        $dbpassword = getenv("DB_PASSWORD");
        $pdo = new \PDO(
            "mysql:dbname=$dbname;host=$dbhost",
            $dbuser,
        $dbpassword,
        );
        $result = $pdo->query("SHOW TABLES");
        $this->assertNotEquals(false, $result);
    }
}