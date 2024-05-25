<?php
namespace Kapigen;
class Configuration {
    public function __construct(private string $test)
    {
    }

    function getTest(): string
    {
        return $this->test;
    }

    function setTest(string $test): void
    {
        $this->test = $test;
    }
}